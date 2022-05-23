package main

import (
	"fmt"
	//"io/ioutil"
    "os/exec"
	"bytes"
	"time"
	//"net/http"
	"strings"
	"encoding/json"
	"strconv"
	"math/rand"
	"math"
	"github.com/go-redis/redis"
)

type FlowFormat []map[string]interface{}
type SimpleMapFormat map[string]interface{}
type HostStats map[string]int

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
	var standardDeviation float64=0
	var alpha float64 = 0.08

	var maxAttackHostsRatio float64 = 0.2
	var R float64 = 0
	
	numberOfLoop:=0;
	for {
		numberOfLoop++;
		if numberOfLoop/30>(numberOfLoop-1)/30 {
			fmt.Println("checked", numberOfLoop, "times")
		}
		
		hubbleFlow, _ := client.Get("newpatch").Result()
		var mapFlow FlowFormat
		json.Unmarshal([]byte(hubbleFlow), &mapFlow)
		
		//phase 1
		//lengthBefore := float64(len(mapFlow))
		mapFlow = filterOnlyRequestTraffic(mapFlow)
    
		//req.body is the input json
		var T float64 = float64(len(mapFlow))
	
		//count number of each host
		hostStats := countHostsAppearance(mapFlow)
		numberOfHost := len(hostStats)
		fmt.Println(hostStats)
		
		//for first scrape, we don't know meanT
		if meanT==0 { meanT = T}
		if (R == 0) {
			R = meanT/float64(numberOfHost)
			fmt.Println("R: ", R)
		}
		fmt.Println("meanT:", meanT, " stdDev:", standardDeviation)    
		//Phase 1
		newMeanT := meanT + alpha * (T - meanT)
		newStandardDeviation := math.Sqrt(alpha * math.Pow(T - meanT, 2) + (1 - alpha)*math.Pow(standardDeviation, 2))
		meanT = newMeanT
		standardDeviation = newStandardDeviation
		threshold := meanT + 3 * standardDeviation
	
		//fmt.Println("meanT:", meanT, "stdDev:",standardDeviation,"threshold: ", threshold)
		fmt.Println("meanT:", meanT, "T:", T, "threshold:", threshold)
		haveSuspected:=true
		if (T <= threshold) {
			//
			R = meanT / float64(numberOfHost)
			fmt.Println("ok")
			haveSuspected = false
			//return res.send('ok')
		}
	
		//Phase 2
		var suspectedHosts string
		if (haveSuspected) {
			fmt.Println("possible attack detected")
			//
			//m
			attackTrafficsRatio := (T - (1-maxAttackHostsRatio) * float64(numberOfHost) * R) / (maxAttackHostsRatio * float64(numberOfHost) * R)
			//mR
			minAttackTraffics := int(attackTrafficsRatio * R)
			fmt.Println("minAttackTraffic:", minAttackTraffics)
			suspectedHosts = getSuspectedHosts(hostStats, minAttackTraffics)    
			client.Set("suspected", suspectedHosts, -1)
		}
		
		time.Sleep(time.Second * 3)
	}
}
func getSuspectedHosts(hostStats HostStats, minAttackTraffics int) (listSuspectedHosts string) {
	listSuspectedHosts = ""
	var buffer bytes.Buffer

	for host, freq := range hostStats {
        if (freq >= minAttackTraffics) {
			buffer.WriteString(host+",")
			//listSuspectedHosts=listSuspectedHosts + host + ","
        }
    }
	listSuspectedHosts = buffer.String()
	lenListSuspectedHosts := len(listSuspectedHosts)
	if lenListSuspectedHosts > 0 {
		listSuspectedHosts = listSuspectedHosts[:lenListSuspectedHosts-1]
	}
	return listSuspectedHosts
}
func countHostsAppearance(mapFlow FlowFormat) (result HostStats) {
	result=HostStats{}
	for _, oneFlow := range mapFlow {
		IP := oneFlow["IP"].(map[string]interface{})["source"].(string)
		if count, isInTheMap := result[IP]; isInTheMap {
			result[IP]=count + 1
		} else {
			result[IP]=1
		}
	}
	return result
}

func filterOnlyRequestTraffic(mapFlow FlowFormat) (result FlowFormat) {
	lastIndex := len(mapFlow) - 1
	for index, oneFlow := range mapFlow {
		if oneFlow["is_reply"]==true {
			mapFlow[index]=mapFlow[lastIndex]
			lastIndex=lastIndex-1	
		}
	}
	return mapFlow[:(lastIndex+1)]
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