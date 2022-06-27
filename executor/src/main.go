package main

import (
	"fmt"
	//"io/ioutil"
    "os/exec"
	//"bytes"
	"encoding/json"
	"time"
	//"net/http"
	"strings"
	"strconv"
	"math/rand"
	"github.com/go-redis/redis"
)
var MAX_CNP_TIME_TO_LIVE=300
type FlowsFormat []FlowFormat
type FlowFormat map[string]interface{}
//type FlowFormat []map[string]interface{}
//key: CNP name
//value: time lived
type WatchingCCNPs map[string]int
func main() {
	fmt.Println("Working")
	rand.Seed(time.Now().UnixNano())

	var (
    	redisUrl     = "redis.kube-system.svc.cluster.local:6379"
    	password = ""
    )
	
    client := redis.NewClient(&redis.Options{
    	Addr:     redisUrl,
    	Password: password,
    	DB:       0,
    })

    _, err := client.Ping().Result()
    if err != nil {
    	fmt.Println(err)
    }

	var listOfWatchingCcnp = WatchingCCNPs{}
	numberOfLoop:=0;
	for {
		numberOfLoop++;
		if numberOfLoop/30>(numberOfLoop-1)/30 {
			fmt.Println("checked", numberOfLoop, "times")
		}
		if len(listOfWatchingCcnp) > 0 {
			//Check if old cnp exists, delete it
			for ccnpName, seconds := range listOfWatchingCcnp {
				if seconds > MAX_CNP_TIME_TO_LIVE {
					delete(listOfWatchingCcnp, ccnpName)
					deleteCcnp(ccnpName)
				} else {
					listOfWatchingCcnp[ccnpName] = seconds + 2
				}
			}
		}
		
		//detect new cnp
		updateNewCcnpToWatchingList(&listOfWatchingCcnp)

		//get suspected and apply cnp
		haveSuspected := true
		
		suspectedFlowsString, err := client.LPop("suspected").Result()
		if err != nil || suspectedFlowsString == "" {
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
				applyCcnp(flow)
			}
		}
		
		time.Sleep(time.Second * 2)
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
	execBashCommand(command)
}
func applyCcnp(flow FlowFormat) {
	unidentifiedFlow := false
	numberOfCcnp := getNumberOfCcnpInString()
	randValue := strconv.Itoa(100+rand.Intn(900))
	command := `cat <<EOF | kubectl apply -f -
apiVersion: 'cilium.io/v2'
kind: CiliumClusterwideNetworkPolicy
metadata:
  name: "blacklist-rule`+numberOfCcnp+randValue+`"
spec:
  endpointSelector:
    matchLabels:`
	for _, label := range flow["destination"].(map[string]interface{})["labels"].([]interface{}) {
		if label.(string)[:4] == "k8s:" {
			realLabel := label.(string)[4:]
			if !strings.Contains(realLabel, "=") {
				unidentifiedFlow=true
				break
			} else {
				realLabel = strings.Replace(realLabel, "=", ": ", -1)
			}
			command += `
      `+realLabel
			fmt.Println(realLabel)
		} else {
			unidentifiedFlow=true
			break
		}
	}
	command +=`
  ingressDeny:
  - fromEndpoints:
    - matchLabels:`
	for _, label := range flow["source"].(map[string]interface{})["labels"].([]interface{}) {
		if label.(string)[:4] == "k8s:" {
			realLabel := label.(string)[4:]
			if !strings.Contains(realLabel, "=") {
				unidentifiedFlow=true
				break
			} else {
				realLabel = strings.Replace(realLabel, "=", ": ", -1)
			}
			command += `
        `+realLabel
			fmt.Println(realLabel)	  
		} else {
			unidentifiedFlow=true
			break
		}
		if strings.Contains(label.(string), "reserved:") {
			unidentifiedFlow = true
			break
		}
	}
	command +=`
  ingress:
  - fromEntities:
    - "all"
EOF`
	if unidentifiedFlow==true {
		fmt.Println("Detect attacker is from outside of cluster:", flow["IP"].(map[string]interface{})["source"].(string))
		return
	}
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

func execBashCommand(command string) (result string, err int) {
	out, cmderr := exec.Command("bash", "-c", command).CombinedOutput()
	if cmderr != nil {
		fmt.Println(cmderr, ":", string(out))
		result = ""
		err=1
	} else {
		result = string(out)
		if len(result) > 0 && result[0]=='"' {
			result = result[1:][:len(result)-2]
		}
		err=0
	}
	return result, err
}