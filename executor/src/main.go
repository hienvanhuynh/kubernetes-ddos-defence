package main

import (
	"fmt"
	//"io/ioutil"
    "os/exec"
	//"bytes"
	"time"
	//"net/http"
	"strings"
	"strconv"
	"math/rand"
	"github.com/go-redis/redis"
)
var MAX_CNP_TIME_TO_LIVE=300

type FlowFormat []map[string]interface{}
//key: CNP name
//value: time lived
type WatchingCNPs map[string]int
func main() {
	fmt.Println("Working")
	
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

	var listOfWatchingCnp = WatchingCNPs{}
	numberOfLoop:=0;
	for {
		numberOfLoop++;
		if numberOfLoop/30>(numberOfLoop-1)/30 {
			fmt.Println("checked", numberOfLoop, "times")
		}

		//Check if old cnp exists, delete it
		for cnpName, seconds := range listOfWatchingCnp {
			if seconds > MAX_CNP_TIME_TO_LIVE {
				delete(listOfWatchingCnp, cnpName)
				deleteCnp(cnpName)
			} else {
				listOfWatchingCnp[cnpName] = seconds + 3
			}
		}
		//detect new cnp
		updateNewCnpToWatchingList(&listOfWatchingCnp)

		//get suspected and apply cnp
		suspectedString, err := client.Get("suspected").Result()
		client.Del("suspected")
		fmt.Println(suspectedString)
		if err!=nil {
			time.Sleep(time.Second * 3)
			continue
		}
		var suspectedIPs []string
		if len(suspectedString) > 0 {
			suspectedIPs = strings.Split(suspectedString, ",")
		} else {
			time.Sleep(time.Second * 3)
			continue
		}

		for _, IP := range suspectedIPs {
			if IP=="" || IP==" " {
				continue
			}
			fmt.Println("Blocking IP:", IP)
			applyCnp(IP)
		}
		
		time.Sleep(time.Second * 3)
	}
}

func updateNewCnpToWatchingList(listOfWatchingCnp *WatchingCNPs) {
	getCNPCommand := "kubectl get cnp --template '{{range .items}}{{.metadata.name}}{{\"\\n\"}}{{end}}' | grep cidr-rule"
	cnpsString, _ := execBashCommand(getCNPCommand)
	cnps := strings.Split(cnpsString, "\n")
	for _, cnp := range cnps {
		if cnp == "" {
			continue
		}
		if _, ok := (*listOfWatchingCnp)[cnp]; ok {
			//Do nothing
		} else {
			(*listOfWatchingCnp)[cnp] = 0
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
func deleteCnp(cnpName string) {
	command:="kubectl delete cnp "+cnpName
	execBashCommand(command)
}
func applyCnp(IP string) {
	numberOfCnp := getNumberOfCnpInString()
	randValue := strconv.Itoa(100+rand.Intn(900))
	command := `cat <<EOF | kubectl apply -f -
apiVersion: "cilium.io/v2"
kind: CiliumNetworkPolicy
metadata:
  name: "cidr-rule`+numberOfCnp+randValue+`"
spec:
  endpointSelector:
    matchLabels:
      app: myapp
  ingress:
  - fromCIDRSet:
    - cidr: 0.0.0.0/0
      except:
      - `+IP+`/32
EOF`
	out, _ := execBashCommand(command);
	fmt.Println(out)
}
func getNumberOfCnpInString() (numberOfCnp string) {
	numberOfCnp, _ = execBashCommand("kubectl get cnp | wc -l")
	//remove unknown character (looks like a space char but it is not) in the last
	numberOfCnp=numberOfCnp[:len(numberOfCnp)-1]
	if (numberOfCnp[0] < '0' || numberOfCnp[0] > '9') {
		numberOfCnp="0"
	}
	return numberOfCnp
}

func execCommand(command string) (result string, err int) {
    commandTokens := strings.Split(command, " ")
	mainCommand, args := commandTokens[0], commandTokens[1:]
	out, cmderr := exec.Command(mainCommand, args...).CombinedOutput()

	if cmderr != nil {
		fmt.Println(cmderr, ":", string(out))
		fmt.Println("Command tokens:", commandTokens)
		result = ""
		err=1
	} else {
		result = string(out)
		if result[0]=='"' {
			result = result[1:][:len(result)-2]
		}
		err=0
	}
	return result, err
}
func execBashCommand(command string) (result string, err int) {
	out, cmderr := exec.Command("bash", "-c", command).CombinedOutput()
	if cmderr != nil {
		fmt.Println(cmderr, ":", string(out))
		result = ""
		err=1
	} else {
		result = string(out)
		if result[0]=='"' {
			result = result[1:][:len(result)-2]
		}
		err=0
	}
	return result, err
}