package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
)

const (
	version = "1.0"
)

func lookupEnv(key string) string {
	for _, v := range os.Environ() {
		kv := strings.Split(v, "=")
		if kv[0] == key {
			return kv[1]
		}
	}
	return ""
}

type openCmd struct {
	cmd  string
	args []string
}

var m = map[string]openCmd{
	"linux":   {cmd: "xdg-open"},
	"windows": {cmd: "rundll32", args: []string{"url.dll,FileProtocolHandler"}},
	"darwin":  {cmd: "open"},
}

func openURLByBrowser(url string) error {
	cmd, ok := m[runtime.GOOS]
	if !ok {
		return errors.New("could not determine how to open URL by browser in this platform")
	}
	args := append(cmd.args, url)
	return exec.Command(cmd.cmd, args...).Start() // #nosec
}

func showVersion() {
	fmt.Println("mymine version", version)
}

type options struct {
	open    []int  `short:"o" long:"open"    description:"Open specified ticket on a web browser"`
	version []bool `short:"v" long:"version" description:"Show version"`
}

func main() {
	var opts options
	parser := flags.NewParser(&opts, flags.Default)

	parser.Name = "mymine"
	parser.Usage = "[OPTIONS]"

	_, err := parser.Parse()
	if err != nil {
		fmt.Println("failed to parse cmd line argument. exit.")
		os.Exit(1)
	}

	if opts.version != nil {
		showVersion()
		os.Exit(1)
	}

	redmineURL := lookupEnv("REDMINE_URL")
	if redmineURL == "" {
		fmt.Println("REDMINE_URL is not specified. exit.")
		os.Exit(1)
	}

	if opts.open != nil {
		url := redmineURL + "issues/" + strconv.Itoa(opts.open[0])
		if err = openURLByBrowser(url); err != nil {
			fmt.Printf("failed to open URL by browser: %s\n", err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

	redmineAPIKey := lookupEnv("REDMINE_API_KEY")
	if redmineAPIKey == "" {
		fmt.Println("REDMINE_API_KEY is not specified. exit.")
		os.Exit(1)
	}

	request := redmineURL + "issues.json?key=" + redmineAPIKey + "&status_id=open&assigned_to_id=me&limit=100"
	fmt.Printf("request = %s\nfetching information...", request)

	content, err := fetchContent(request)
	if err != nil {
		fmt.Printf("failed to fetch content: %s", err.Error())
		os.Exit(1)
	}

	showContent(content)
}

func fetchContent(url string) (map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("failed to fetch information: %s\n", err.Error())
		os.Exit(1)
	}
	defer func() {
		_ = resp.Body.Close() // #nosec
	}()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err.Error())
	}

	var buf map[string]interface{}
	err = json.Unmarshal(content, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %s", err.Error())
	}
	return buf, nil
}

func showContent(m map[string]interface{}) {
	issues := m["issues"].([]interface{})
	for _, v := range issues {
		issue := v.(map[string]interface{})
		id := int(issue["id"].(float64))
		status := issue["status"].(map[string]interface{})
		fmt.Printf("[#%d] %11s %s\n", id, status["name"], issue["subject"])
	}
}
