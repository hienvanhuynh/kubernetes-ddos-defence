package main

import (
	"fmt"
    "os/exec"
	"time"
	"strings"
	"encoding/json"
	"github.com/go-redis/redis"
)

type FlowFormat []map[string]interface{}
func main() {
	fmt.Println("Working")
	MIN_PATCH_ID := 0
	MAX_PATCH_ID := 100000000
	patchid := MIN_PATCH_ID

	// Create Redis Client
    var (
    	redisUrl     = "kdd-redis.kube-system.svc.cluster.local:6379"
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
		
		for {

			hubbleFlow, getHubbleFlowErr := getHubbleFlow()

			if getHubbleFlowErr != nil {
				time.Sleep(time.Second * 5)
				continue
			}

			var mapFlow FlowFormat
			json.Unmarshal([]byte(hubbleFlow), &mapFlow)
			
			//for _, oneFlow := range mapFlow {
			//	delete(oneFlow, "node_name")
			//	delete(oneFlow, "reply")
			//	delete(oneFlow, "event_type")
			//	delete(oneFlow, "Summary")
			//	delete(oneFlow, "trace_observation_point")
			//}
	
			shortJsonStringFlow, _ := json.Marshal(mapFlow)
			redisError := client.Set("newpatch", shortJsonStringFlow, -1).Err()

			if redisError != nil {
				break
			}
	
			redisError = client.Set("patchid", patchid, -1).Err()
			if redisError != nil {
				break
			}
	
			patchid += 1
			if patchid > MAX_PATCH_ID {
				patchid = MIN_PATCH_ID
			}

			time.Sleep(time.Second * 5)
		}
	}
}

func getHubbleFlow() (result string, returnerr error) {

	//Get hubble relay server
	command := "kubectl -n kube-system get svc hubble-relay -o jsonpath='{.spec.clusterIP}'"
	hubbleRelayIP, getHubbleIPErr := execBashCommand(command)

	if getHubbleIPErr != nil || hubbleRelayIP == "" {
		fmt.Println("ERROR: No hubble relay detected, will not be able to work")
		return "", getHubbleIPErr
	}

	//Get flow
	command = "kubectl exec " + getPodName("kube-system", "k8s-app=cilium") + " -n kube-system -- hubble --server " + hubbleRelayIP + ":80 observe --since 5.5s --verdict FORWARDED -o json"
	rawFlows, getFlowErr := execBashCommand(command)

	if getFlowErr != nil {
		return "", getFlowErr
	}

	splitedFlows := strings.Split(rawFlows, "\n")
	formatedFlows := "[" + strings.Join(splitedFlows[1:][:len(splitedFlows)-2], ",\n")+"]"

	return formatedFlows, nil
}

func getPodName(namespace string, labels ...string) (result string) {

	command := "kubectl get pod -n "+namespace+" -o jsonpath=\"{.items[0].metadata.name}\""
	for _, label := range labels {
		command = command + " -l " + label
	}

	result, _ = execBashCommand(command)
	return result
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
	return result, cmderr
}