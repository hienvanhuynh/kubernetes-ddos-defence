package main

import (
    "fmt"
    "net/http"
	"strings"
	"os/exec"
)
  
type WatchingCNPs map[string]int
func main() {
    fmt.Println("Listen on Port 7080")
    http.HandleFunc("/", homeGUI)
	http.HandleFunc("/blocklist", getBlockList)
    http.ListenAndServe(":7080", nil)
}

func getBlockList(w http.ResponseWriter, r *http.Request) {
	cnpList := getCnpList()
	response := ""
	for _, cnp := range cnpList {
		if len(response)!=0 {
			response+="\n"
		}
		response += cnp
	}

	fmt.Fprintf(w, response)
}
func homeGUI(w http.ResponseWriter, r *http.Request) {

	htmlResponse := `
	<!DOCTYPE html>
<html>
<head>

<head>
<body>
    <h1>Kubernetes DDOS Defence UI</h1>
    <p>UI version: 1.0.0</p>
    Blocking IPs:
    <ul>
	`
	cnpList := getCnpList()

	for _, cnp := range cnpList {
		htmlResponse+="<li>"+cnp+"</li>"
	}

	htmlResponse+=`    </ul>
	</body>
	</html>`

    fmt.Fprintf(w, htmlResponse)
}

func getCnpList() (cnpList []string) {
	command := "kubectl -n kube-system get cnp -o jsonpath='{range .items[*]}{.metadata.name}{\"\\n\"}{end}'"
	cnpsString, _ := execBashCommand(command)
	cnpList = strings.Split(cnpsString, "\n")
	return cnpList
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