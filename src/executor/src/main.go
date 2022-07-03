package main

import (
	"fmt"
    "os/exec"
	"encoding/json"
	"time"
	"strings"
	"strconv"
	"math/rand"
	"io"
	"github.com/go-redis/redis"
)

var MAX_CNP_TIME_TO_LIVE=600
type FlowsFormat []FlowFormat
type FlowFormat map[string]interface{}

//key: CNP name
//value: time lived
type WatchingCCNPs map[string]int

func main() {
	rand.Seed(time.Now().UnixNano())

	var (
    	redisUrl     = "redis.kube-system.svc.cluster.local:6379"
    	password = ""
    )

	for {
		client := redis.NewClient(&redis.Options{
			Addr:     redisUrl,
			Password: password,
			DB:       0,
		})
	
		_, err := client.Ping().Result()
		if err != nil {
			fmt.Println("ERROR: Could not found redis service, will not be able to work")
			fmt.Println(err)
			//Wait for sometime before retry, if don't do this then the log will be spammed
			time.Sleep(time.Second * 5)
			continue
		}

		fmt.Println("Working")

		var listOfWatchingCcnp = WatchingCCNPs{}

		for {
			time.Sleep(time.Second * 3)
			if len(listOfWatchingCcnp) > 0 {
				//Check if old cnp exists, delete it
				for ccnpName, seconds := range listOfWatchingCcnp {
					if seconds > MAX_CNP_TIME_TO_LIVE {
						delete(listOfWatchingCcnp, ccnpName)
						deleteCcnp(ccnpName)
					} else {
						listOfWatchingCcnp[ccnpName] = seconds + 3
					}
				}
			}
			
			//detect new cnp
			updateNewCcnpToWatchingList(&listOfWatchingCcnp)
	
			//get suspected and apply cnp
			haveSuspected := true
			
			suspectedFlowsString, err := client.LPop("suspected").Result()
			if err != nil {
				if err == io.EOF {
					break
				}
			} else if suspectedFlowsString == "" {
				haveSuspected = false
			}
	
			var suspectedFlows FlowsFormat
			json.Unmarshal([]byte(suspectedFlowsString), &suspectedFlows)
			
			//suspectedString, err := client.Get("suspected").Result()
			//fmt.Println(suspectedString)
			if len(suspectedFlows) <= 0 {
				haveSuspected = false
			}
	
			if haveSuspected == true {
				for _, flow := range suspectedFlows {
					fmt.Println("Blocking flow:", flow)
					applyCcnp(flow, listOfWatchingCcnp)
				}
			}
		}
	}
}

func updateNewCcnpToWatchingList(listOfWatchingCcnp *WatchingCCNPs) {	
	getCCNPCommand := "kubectl get ccnp --template '{{range .items}}{{.metadata.name}}{{\"\\n\"}}{{end}}' | grep blacklist-rule"
	ccnpsString, _ := execBashCommand(getCCNPCommand)
	ccnps := strings.Split(ccnpsString, "\n")
	for _, ccnp := range ccnps {
		if ccnp == "" {
			continue
		}
		if _, ok := (*listOfWatchingCcnp)[ccnp]; ok {
			//Do nothing
		} else {
			(*listOfWatchingCcnp)[ccnp] = 0
		}
	}
}

func isSeparatorInIPList(sep rune) (result bool) {
	if sep=='[' || sep ==']' || sep=='"' {
		return true
	} else {
		return false
	}
}
func deleteCcnp(ccnpName string) {
	command:="kubectl delete ccnp "+ccnpName
	fmt.Println("delete command:"+command)
	execBashCommand(command)
}

func getPolicy(policyName string) (policy string){
	command := "kubectl get ccnp " + policyName + " -o yaml"
	out, err := execBashCommand(command)
	if err != nil {
		fmt.Println("Failed to get policy content")
	}
	return out
}

func testIfThisPolicyAlreadyExists(thisSpec string, listOfWatchingCcnp WatchingCCNPs) (exists bool) {
	for policyName,_ := range listOfWatchingCcnp {
		watchingPolicy := getPolicy(policyName)
		
		if strings.Contains(watchingPolicy, thisSpec) {
			return true
		}
	}
	return false
}

func applyCcnp(flow FlowFormat, listOfWatchingCcnp WatchingCCNPs) {
	unidentifiedFlow := false
	worldFlow := false
	numberOfCcnp := getNumberOfCcnpInString()
	randValue := strconv.Itoa(100+rand.Intn(900))

	commandHeader := `cat <<EOF | kubectl apply -f -
apiVersion: 'cilium.io/v2'
kind: CiliumClusterwideNetworkPolicy
metadata:
  name: "blacklist-rule`+numberOfCcnp+randValue+`"`
    policySpec :=`
spec:
  endpointSelector:
    matchLabels:`
	for _, label := range flow["destination"].(map[string]interface{})["labels"].([]interface{}) {
		realLabel := label.(string)[4:]

		if label.(string)[:4] == "k8s:" {
			if !strings.Contains(realLabel, "=") {
				unidentifiedFlow=true
			} else {
				realLabel = strings.Replace(realLabel, "=", ": ", -1)
			}

		} else {
			realLabel = strings.Replace(label.(string), "=", ": ", -1)
		}
		policySpec += `
      `+realLabel
	}
	policySpec +=`
  ingress:
  - fromEntities:
    - all`
	policySpec +=`
  ingressDeny:
  - fromEndpoints:
    - matchLabels:`
	for _, label := range flow["source"].(map[string]interface{})["labels"].([]interface{}) {
		if label.(string)[:4] == "k8s:" {
			realLabel := label.(string)[4:]
			if !strings.Contains(realLabel, "=") {
				unidentifiedFlow=true
			} else {
				realLabel = strings.Replace(realLabel, "=", ": ", -1)
			}
			policySpec += `
        `+realLabel
			fmt.Println(realLabel)	  
		}
		if strings.Contains(label.(string), "reserved:world") {
			worldFlow = true
			break
		}
	}
	
	if worldFlow == true {
		unidentifiedFlow = false
		policySpec = `
spec:
  endpointSelector:
    matchLabels:`
		for _, label := range flow["destination"].(map[string]interface{})["labels"].([]interface{}) {
			if label.(string)[:4] == "k8s:" {
				realLabel := label.(string)[4:]
				if !strings.Contains(realLabel, "=") {
					continue
				} else {
					realLabel = strings.Replace(realLabel, "=", ": ", -1)
				}
				policySpec += `
      `+ realLabel
				fmt.Println(realLabel)
			}
		}
		policySpec +=`
  ingress:
  - fromEntities:
    - all`
		policySpec +=`
  ingressDeny:
  - fromCIDR:`
		externalIP := flow["IP"].(map[string]interface{})["source"].(string)
		policySpec += `
    - ` + externalIP + `/32`
	}
	
	if unidentifiedFlow == true {
		fmt.Println("attack flow is unknown")
		return
	}

	exists := testIfThisPolicyAlreadyExists(policySpec, listOfWatchingCcnp)
	if exists == true {
		fmt.Println("Policy is already existed")
		return;
	}
	
	command := commandHeader + policySpec + `
EOF`

	out, _ := execBashCommand(command);
	fmt.Println(out)
}
func getNumberOfCcnpInString() (numberOfCcnp string) {
	numberOfCcnp, _ = execBashCommand("kubectl get ccnp | wc -l")
	//remove unknown character (looks like a space char but it is not) in the last
	numberOfCcnp=numberOfCcnp[:len(numberOfCcnp)-1]
	if (numberOfCcnp[0] < '0' || numberOfCcnp[0] > '9') {
		numberOfCcnp="0"
	}
	return numberOfCcnp
}

func execBashCommand(command string) (result string, err error) {
	out, cmderr := exec.Command("bash", "-c", command).CombinedOutput()
	if cmderr != nil {
		//fmt.Println(cmderr, ":", string(out))
		result = ""
	} else {
		result = string(out)
		if len(result) > 0 && result[0]=='"' {
			result = result[1:][:len(result)-2]
		}
	}
	return result, cmderr
}