package main

import (
	"fmt"
	//"io/ioutil"
    "os/exec"
	//"bytes"
	"time"
	"reflect"
	//"net/http"
	"strings"
	"encoding/json"
	"strconv"
	"math/rand"
	"math"
	//"os"
	"github.com/go-redis/redis"
)

type FlowsFormat []FlowFormat
type FlowFormat map[string]interface{}
//index is which flow this is pointing to, value is the counted appearance number
type FlowsStats map[int]int

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
	var meanT float64 = 0
	var standardDeviation float64 = 0
	var alpha float64 = 0.08

	//In case stdDev is low, we use this number to keep it not too low
	var standardDeviationBias float64 = 5
	//Every user may as the same time increase 5 access, then we tolerate them
	var tolerationTrafficIncreasementBias float64 = 5
	//var numberOfDnsService = 1

	var maxAttackHostsRatio float64 = 0.1
	//Usual traffic of a client
	var R float64 = 0
	
	savedPatchId := "-1"
	numberOfLoop := 0;
	for {
		numberOfLoop++;
		if numberOfLoop/30>(numberOfLoop-1)/30 {
			fmt.Println("checked", numberOfLoop, "times")
		}
		
		hubbleFlows, _ := client.Get("newpatch").Result()
		patchid, _ := client.Get("patchid").Result()

		if savedPatchId != patchid {
			savedPatchId = patchid
			//Start to analyze
			var mapFlows FlowsFormat
			json.Unmarshal([]byte(hubbleFlows), &mapFlows)
			
			//phase 1
			//lengthBefore := float64(len(mapFlow))
			mapFlows = filterOnlyRequestTraffic(mapFlows)
		
			//req.body is the input json
			var T float64 = float64(len(mapFlows))
		
			//count number of each host
			flowsStats := countHostsAppearance(mapFlows)
			numberOfFlow := len(flowsStats)
			fmt.Println(flowsStats)
			
			//for first scrape, we don't know meanT
			if meanT==0 { meanT = T}
			if (R == 0) {
				R = meanT/float64(numberOfFlow)
				fmt.Println("R: ", R)
			}
			//fmt.Println("meanT:", meanT, " stdDev:", standardDeviation)    
			//Phase 1 identify there are unexpected hight traffic
			//save StandardDeviation for phase 1.1
			tolerationTrafficIncreasement := standardDeviation / float64(numberOfFlow)
			if tolerationTrafficIncreasement < tolerationTrafficIncreasementBias {
				tolerationTrafficIncreasement = tolerationTrafficIncreasementBias
			}
			newMeanT := meanT + alpha * (T - meanT)
			newStandardDeviation := math.Sqrt(alpha * math.Pow(T - meanT, 2) + (1 - alpha)*math.Pow(standardDeviation, 2))
			threshold := newMeanT + 3 * (newStandardDeviation + standardDeviationBias)
		
			//fmt.Println("meanT:", meanT, "stdDev:",standardDeviation,"threshold: ", threshold)
			//fmt.Println("meanT:", meanT, "T:", T, "threshold:", threshold)
			haveSuspected:=true
			//check if T is not exceed the threshold and also check if the increasing of traffic is not purely caused by increase number of clients
			if (T <= threshold || T <= float64(numberOfFlow) * (R + tolerationTrafficIncreasement)) {
				//
				meanT = newMeanT
				standardDeviation = newStandardDeviation
				R = meanT / float64(numberOfFlow)
				//fmt.Println("ok")
				haveSuspected = false
			}
	
			//Phase 2 filter the highest possible attack flow
			var suspectedFlows FlowsFormat
			if (haveSuspected) {
				fmt.Println("possible attack detected")
				//A
				numberOfAttackFlow := maxAttackHostsRatio * float64(numberOfFlow)
				if numberOfAttackFlow < 1 { numberOfAttackFlow = 1 }
				//Each traffic caused by attacker will create a similar dns flow
				//So we must add this to make the argorithm works correctly
				numberOfAttackFlow = numberOfAttackFlow * 2
				//mR
				minAttackTraffics := int((T - (float64(numberOfFlow) - numberOfAttackFlow) * R) / numberOfAttackFlow)
				fmt.Println("minAttackTraffic:", minAttackTraffics)
				suspectedFlows = getSuspectedFlows(mapFlows, flowsStats, minAttackTraffics)
				
				if len(suspectedFlows) == 0 {
					meanT = newMeanT
					standardDeviation = newStandardDeviation
					R = meanT / float64(numberOfFlow)	
				} else {
					fmt.Println("Attack confirmed")
					suspectedFlowsJson, _ := json.Marshal(suspectedFlows)
					fmt.Println(string(suspectedFlowsJson))
					client.RPush("suspected", string(suspectedFlowsJson))
				}
			}
		}
		
		time.Sleep(time.Second * 3)
	}
}

func (flow1 FlowFormat) Equals(flow2 FlowFormat) bool {
	if reflect.DeepEqual(flow1["IP"], flow2["IP"]) && 
		reflect.DeepEqual(flow1["source"].(map[string]interface{})["labels"], flow2["source"].(map[string]interface{})["labels"]) {
			return true
	}
	return false
}

func getSuspectedFlows(flows FlowsFormat, flowsStats FlowsStats, minAttackTraffics int) (listSuspectedFlows FlowsFormat) {
	listSuspectedFlows=FlowsFormat{}
	//var buffer bytes.Buffer
	for flowIndex, freq := range flowsStats {
        if (freq >= minAttackTraffics) {
			listSuspectedFlows = append(listSuspectedFlows, flows[flowIndex])
			//buffer.WriteString(host+",")
			//listSuspectedHosts=listSuspectedHosts + host + ","
        }
    }
	//listSuspectedHosts = buffer.String()
	//lenListSuspectedFlows := len(listSuspectedFlows)
	//To delete the last comma
	//if lenListSuspectedHosts > 0 {
	//	listSuspectedHosts = listSuspectedHosts[:lenListSuspectedHosts-1]
	//}
	return listSuspectedFlows
}
func countHostsAppearance(mapFlows FlowsFormat) (result FlowsStats) {
	result=FlowsStats{}
	for oneFlowIndex, oneFlow := range mapFlows {
		found := false
		for indexFlow, numberAppearance := range result {
			if mapFlows[indexFlow].Equals(oneFlow) {
				result[indexFlow] = numberAppearance+1
				found = true
				break
			}
		}
		if found==true {
			continue
		}
		result[oneFlowIndex] = 1
	}
	return result
}

func filterOnlyRequestTraffic(mapFlows FlowsFormat) (result FlowsFormat) {
	lastIndex := len(mapFlows) - 1
	for index, oneFlow := range mapFlows {
		if oneFlow["is_reply"]==true {
			mapFlows[index]=mapFlows[lastIndex]
			lastIndex=lastIndex-1	
		} else if oneFlow["verdict"]=="DROPPED" {
			mapFlows[index]=mapFlows[lastIndex]
			lastIndex=lastIndex-1
		}
	}
	return mapFlows[:(lastIndex+1)]
}

func isSeparatorInIPList(sep rune) (result bool) {
	if sep=='[' || sep ==']' || sep=='"' {
		return true
	} else {
		return false
	}
}
func applyCnp(IP string) {
	command := `cat <<EOF | kubectl apply -f -
apiVersions: "cilium.io/v2"
kind: CiliumNetworkPolicy
metadata:
  name: "cidr-rule`+strconv.Itoa(getNumberOfCnp()+rand.Intn(1000))+`"
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
func getNumberOfCnp() (result int) {
	numberString, _ := execCommand("kubectl get cnp | wc -l")
	numberOfCnp, _ := strconv.Atoi(numberString)

	return numberOfCnp
}
func getHubbleFlow() (result string) {
	command := "kubectl exec " + getPodName("kube-system", "k8s-app=cilium") + " -n kube-system -- hubble observe --since 3s -o json"
	rawFlows, _ := execCommand(command)
	splitedFlow := strings.Split(rawFlows, "\n")
	
	var formatedFlow = "[" + strings.Join(splitedFlow[1:][:len(splitedFlow)-2], ",\n")+"]"
	return formatedFlow
}

func getPodName(namespace string, labels ...string) (result string) {
	command := "kubectl get pod -n "+namespace+" -o jsonpath=\"{.items[0].metadata.name}\""
	for _, label := range labels {
		command = command + " -l " + label
	}

	result, _ = execCommand(command)
	return result
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