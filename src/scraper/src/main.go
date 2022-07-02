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
    /*var (
    	host     = getEnv("REDIS_HOST", "http://redis.kube-system.pod.cluster.local:6379")
    	port     = string(getEnv("REDIS_PORT", "6379"))
    	password = getEnv("REDIS_PASSWORD", "")
    )*/
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

	numberOfLoop:=0;
	for {
		numberOfLoop++;
		//if numberOfLoop/30>(numberOfLoop-1)/30 {
		//	fmt.Println("checked", numberOfLoop, "times")
		//}
		hubbleFlow := getHubbleFlow()
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
		client.Set("newpatch", shortJsonStringFlow, -1)
		client.Set("patchid", patchid, -1)
		patchid += 1
		if patchid > MAX_PATCH_ID {
			patchid = MIN_PATCH_ID
		}
		time.Sleep(time.Second * 5)
	}
}

func getHubbleFlow() (result string) {
	//Get hubble relay server
	command := "kubectl -n kube-system get svc hubble-relay -o jsonpath='{.spec.clusterIP}'"
	hubbleRelayIP, err := execBashCommand(command)
	if err != 0 || hubbleRelayIP == "" {
		fmt.Println("WARNING: No hubble relay detected, may not work correctly")
	}
	//Get flow
	command = "kubectl exec " + getPodName("kube-system", "k8s-app=cilium") + " -n kube-system -- hubble --server " + hubbleRelayIP + ":80 observe --since 5.5s --verdict FORWARDED -o json"
	rawFlows, _ := execBashCommand(command)
	splitedFlow := strings.Split(rawFlows, "\n")
	
	var formatedFlow = "[" + strings.Join(splitedFlow[1:][:len(splitedFlow)-2], ",\n")+"]"
	return formatedFlow
}

func getPodName(namespace string, labels ...string) (result string) {
	command := "kubectl get pod -n "+namespace+" -o jsonpath=\"{.items[0].metadata.name}\""
	for _, label := range labels {
		command = command + " -l " + label
	}

	result, _ = execBashCommand(command)
	return result
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