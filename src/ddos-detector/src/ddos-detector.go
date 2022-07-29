package main

import (
	"fmt"
    "os/exec"
	"time"
	"reflect"
	"strings"
	"encoding/json"
	"strconv"
	"math/rand"
	"math"
	"github.com/go-redis/redis"
)

type FlowsFormat []FlowFormat
type FlowFormat map[string]interface{}
//index is which flow this is pointing to, value is the counted appearance number
type FlowsStats map[int]int
type Candidate struct {
	flow FlowFormat
	numberOfTimes int
}
type Candidates []Candidate


func main() {
	fmt.Println("Working")

    var (
    	redisUrl     = "kdd-redis.kube-system.svc.cluster.local:6379"
    	password = ""
    )
	
	var meanT float64 = 0
	var standardDeviation float64 = 0
	var alpha float64 = 0.15

	//In case stdDev is low, we use this number to keep it not too low
	//var standardDeviationBias float64 = 2
	//Every user may as the same time increase 5 access, then we tolerate them
	//var tolerationTrafficIncreasementBias float64 = 5
	//var numberOfDnsService = 1

	var maxAttackHostsRatio float64 = 0.2
	//Usual traffic of a client
	var R float64 = 0

	var candidates Candidates

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
			time.Sleep(time.Second * 3)
			continue
		}

		savedPatchId := "-1"
		numberOfLoop := 0;
		patchid := "-1"
		for {
			numberOfLoop++;
			if numberOfLoop/30>(numberOfLoop-1)/30 {
				fmt.Println("checked", numberOfLoop, "times")
			}

			hubbleFlows, redisError := client.Get("newpatch").Result()
			if redisError != nil {
				fmt.Println("Error when get newpatch: ", redisError)
				break
			}
	
			patchid, redisError = client.Get("patchid").Result()
			fmt.Println("patch: ", patchid)
			
			if redisError != nil {
				fmt.Println("Error when get patchid: ", redisError)
				break
			}

			var mapFlows FlowsFormat
			json.Unmarshal([]byte(hubbleFlows), &mapFlows)
			
			//phase 1
			//lengthBefore := float64(len(mapFlow))
			//mapFlows = filterMainTraffic(mapFlows)
			mapFlows = filterMainTraffic(mapFlows)

			var T float64 = float64(len(mapFlows))
			if savedPatchId == patchid || T == 0 {
				fmt.Println("No new patch")
			} else {
				savedPatchId = patchid
				//Start to analyze
				
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
				//Phase 1 identify there are unexpected hight traffic
				//save StandardDeviation for phase 1.1
				//tolerationTrafficIncreasement := standardDeviation / float64(numberOfFlow)
				//if tolerationTrafficIncreasement < tolerationTrafficIncreasementBias {
				//	tolerationTrafficIncreasement = tolerationTrafficIncreasementBias
				//}
			
				newMeanT := meanT + alpha * (T - meanT)
				newStandardDeviation := math.Sqrt(alpha * math.Pow(T - meanT, 2) + (1 - alpha)*math.Pow(standardDeviation, 2))
				//threshold := newMeanT + 3 * newStandardDeviation
				threshold := meanT + 3 * standardDeviation
			
				//fmt.Println("meanT:", meanT, "stdDev:",standardDeviation,"threshold: ", threshold)
				fmt.Println("meanT:", meanT, "T:", T, "threshold:", threshold)
				haveSuspected:=true
				//check if T is not exceed the threshold
				// Possibly and also check if the increasing of traffic is not purely caused by increase number of clients
				// || T <= float64(numberOfFlow) * (R + tolerationTrafficIncreasement)
				if T <= threshold {
					//
					meanT = newMeanT
					standardDeviation = newStandardDeviation
					R = meanT / float64(numberOfFlow)
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
					suspectedFlows, candidates = getSuspectedFlows(&mapFlows, &flowsStats, minAttackTraffics, &candidates)
					
					if len(suspectedFlows) == 0 && len(candidates) == 0 {
						meanT = newMeanT
						standardDeviation = newStandardDeviation
						R = meanT / float64(numberOfFlow)	
					} else {
						fmt.Println("Attack confirmed")
						suspectedFlowsJson, _ := json.Marshal(suspectedFlows)
						fmt.Println(string(suspectedFlowsJson))
						redisError = client.RPush("suspected", string(suspectedFlowsJson)).Err()
						if redisError != nil {
							break
						}
					}
				}
			}
			
			time.Sleep(time.Second * 3)
		}
	}
}

func (flow1 FlowFormat) Equals(flow2 FlowFormat) bool {
	if reflect.DeepEqual(flow1["IP"], flow2["IP"]) && 
		reflect.DeepEqual(flow1["source"].(map[string]interface{})["labels"], flow2["source"].(map[string]interface{})["labels"]) {
			return true
	}
	return false
}

func getSuspectedFlows(flows *FlowsFormat, flowsStats *FlowsStats, minAttackTraffics int, oldCandidates *Candidates) (listSuspectedFlows FlowsFormat, newCandidates Candidates) {
	listSuspectedFlows=FlowsFormat{}
    
	lenOldCandidates := len(*oldCandidates)

	for flowIndex, freq := range *flowsStats {
        
		if (freq >= minAttackTraffics) {
			
			newCandidateFlow := (*flows)[flowIndex]

			numAppear := 1

			for index := 0; index < lenOldCandidates; index++ {
				
				if newCandidateFlow.Equals((*oldCandidates)[index].flow) {
					
					numAppear = numAppear + (*oldCandidates)[index].numberOfTimes
					
					lenOldCandidates--
					(*oldCandidates)[index] = (*oldCandidates)[lenOldCandidates]
					
					break
				}
			}

			//If this flow is the candidate 3 times then decide this is a suspected flow
			if numAppear >= 3 {

				listSuspectedFlows = append(listSuspectedFlows, newCandidateFlow)
		
			} else {

				//If flow is the candidate less than 3 times then add it to new candidates list
				newCandidate := Candidate{newCandidateFlow, numAppear}
				newCandidates = append(newCandidates, newCandidate)
			}
		}
    }
	
	return listSuspectedFlows, newCandidates
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

func filterMainTraffic(mapFlows FlowsFormat) (result FlowsFormat) {
	lastIndex := len(mapFlows) - 1
	for index, oneFlow := range mapFlows {
		if oneFlow["is_reply"]==true {
			mapFlows[index]=mapFlows[lastIndex]
			lastIndex=lastIndex-1	
		} else {
			is_ingress := false
			for _, label := range oneFlow["source"].(map[string]interface{})["labels"].([]interface{}) {
				if label.(string) == "k8s:app.kubernetes.io/instance=ingress-nginx" || label.(string) == "k8s:app.kubernetes.io/name=ingress-nginx" {
						is_ingress = true
						break
					}
				if strings.Contains(label.(string), "ingress-nginx") {
					is_ingress = true
					break
				}
			}
			if is_ingress == true {
				mapFlows[index]=mapFlows[lastIndex]
				lastIndex=lastIndex-1
			}
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
func execBashCommand(command string) (result string, err error) {
	out, cmderr := exec.Command("bash", "-c", command).CombinedOutput()
	if cmderr != nil {
		fmt.Println(cmderr, ":", string(out))
		result = ""
	} else {
		result = string(out)
		if result[0]=='"' {
			result = result[1:][:len(result)-2]
		}
	}
	return result, err
}