package main

import (
    "fmt"
    "net/http"
	"strings"
	"os/exec"
	"encoding/json"
)

type PolicyConfig map[string]interface{}
func main() {
    fmt.Println("Listen on Port 7080")
    http.HandleFunc("/", homeGUI)
	http.HandleFunc("/blocklist", getBlockList)
    http.ListenAndServe(":7080", nil)
}

func getBlockList(w http.ResponseWriter, r *http.Request) {
	policyList := getCcnpList()
	fmt.Println("blocked flows:", policyList)
	
	response, _ := json.Marshal(policyList)

	fmt.Fprintf(w, string(response))
}
func homeGUI(w http.ResponseWriter, r *http.Request) {

	htmlResponse := `
	<!DOCTYPE html>
	<html>
	<head>
	
	<style>
	body {
		margin:0;
	}
	*, *:before, *:after {
		box-sizing: border-box;
	}
	.wrapper {
		border: 5px;
	}
	li {
		background-color: white;
		padding:10px;
	}
	li:nth-child(odd) {
		background-color: #d0e9f7;
	}
	li:nth-child(even) {
		background-color: #FFFFFF;
	}
	.light-sep {
	  border: 2px solid #EEEEEE;
	  border-radius: 2px;
	  height: 2px;
	  width:90%%;
	  left:5%%;
	  position: relative;
	  top: 20px;
	}
	</style>
	<head>
	<body style="overflow:hidden;">
		
		<div style="position: absolute; width:100%%; height: 10%%; background-color:#BCF; box-shadow:0px 0px 3px 3px #BCF; display: flex; justify-content: left; align-items: center;">
		<h1>Kubernetes DDOS Defence UI</h1>
		</div>
		<div style="position: absolute; width:99%%; height:90%%; left:0.5%%; top:11%%; border-radius:10px; background-color: #FFFFFF; box-shadow: 0px 0px 5px 5px #DDDDDD; overflow: hidden;">
	
		<div style="position: relative; left:10px; top:10px; height:20px; display: flex; align-items: center;">Blocking Flows:
		</div>
		<div class="light-sep" style=""></div>
		<div style="position: relative; top:20px; height:90%%; overflow: hidden; overflow-y: scroll; background-color: #FFFFFF; margin: 0; padding: 0;">
		<ul style="list-style-type: none; margin: 0; padding: 0;">
	`
	policyList := getCcnpList()
	fmt.Println(policyList)
	for id, policy := range policyList {
		fmt.Println(id, policy)
		policyJson, _ := json.Marshal(policy)

		htmlResponse+="<li>"+string(policyJson)+"</li>"
	}

	htmlResponse+=`        </ul>
    </div>
    </div>
</body>
</html>`

    fmt.Fprintf(w, htmlResponse)
}

func getCcnpList() (policyListJson []map[string]interface{}) {
	command := "kubectl get ccnp -o jsonpath='{range .items[*]}{.metadata.annotations.kubectl\\.kubernetes\\.io/last-applied-configuration}{end}'"
	policiesString, _ := execBashCommand(command)
	policyListOpenJsonStr := "[" + strings.Replace(policiesString, "\n", ",", -1)
	policyListCloseJsonStr := policyListOpenJsonStr[:len(policyListOpenJsonStr)-1] + "]"
	json.Unmarshal([]byte(policyListCloseJsonStr), &policyListJson)
	return policyListJson
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