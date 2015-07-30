package main

import (
	"fmt"
	"menteslibres.net/gosexy/rest"
	"os"
	"strings"
)

func getEnvVar(varName string) (result string) {
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if pair[0] == varName {
			return pair[1]
		}
	}
	return ""
}

func main() {
	redmineUrl := getEnvVar("REDMINE_URL")
	if redmineUrl == "" {
		fmt.Println("REDMINE_URL is not specified.")
		return
	}

	redmineApiKey := getEnvVar("REDMINE_API_KEY")
	if redmineApiKey == "" {
		fmt.Println("REDMINE_API_KEY is not specified.")
		return
	}

	request := redmineUrl + "issues.json?key=" + redmineApiKey + "&status_id=open&assigned_to_id=me&limit=100"
	fmt.Println("request = ", request)
	var buf map[string]interface{}
	rest.Get(&buf, request, nil)

	issues := buf["issues"].([]interface{})
	for _, v := range issues {
		issue := v.(map[string]interface{})
		fmt.Println(issue["id"], issue["subject"])
	}
}
