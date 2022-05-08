package main

import (
	"fmt"
	"io/ioutil"
    "os/exec"
	"bytes"
	"time"
	"net/http"
	"strings"
	"encoding/json"
	"strconv"
	"math/rand"
)

type FlowFormat []map[string]interface{}
func main() {
	fmt.Println("Working...")
	detectionUrl := "http://ddos-detection.kube-system.svc.cluster.local:5060/newpatch"

	numberOfLoop:=0;
	for {
		numberOfLoop++;
		if numberOfLoop/30>(numberOfLoop-1)/30 {
			fmt.Println("checked", numberOfLoop, "times")
		}
		hubbleFlow := getHubbleFlow()
		var mapFlow FlowFormat
		json.Unmarshal([]byte(hubbleFlow), &mapFlow)
		
		for _, oneFlow := range mapFlow {
			delete(oneFlow, "source")
			delete(oneFlow, "destination")
			delete(oneFlow, "node_name")
			delete(oneFlow, "reply")
			delete(oneFlow, "event_type")
			delete(oneFlow, "Summary")
			delete(oneFlow, "trace_observation_point")
			delete(oneFlow, "verdict")
		}

		shortJsonStringFlow, _ := json.Marshal(mapFlow)

		res, err := http.Post(detectionUrl, "application/json", bytes.NewBuffer(shortJsonStringFlow))
		if err != nil {
			fmt.Println("Post request error:", err)
			continue
		}
		
		resContent,err := ioutil.ReadAll(res.Body)
        if err != nil {
    		fmt.Println("Parse response body error:", err)
			continue
		}
		
		if (string(resContent)=="ok") {
			continue
		}

		fmt.Println("Detected:", string(resContent))
		IPs := strings.FieldsFunc(string(resContent), isSeparatorInIPList)
		for _, IP := range IPs {
			if IP=="" || IP==" " {
				continue
			}
			fmt.Println("Blocking IP:", IP)
			applyCnp(IP)
		}
		
		time.Sleep(time.Second * 3)
	}
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
apiVersion: "cilium.io/v2"
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