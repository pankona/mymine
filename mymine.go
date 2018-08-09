package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
	"menteslibres.net/gosexy/rest"
	"github.com/pkg/errors"
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

func openURLByBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = errors.New("cannot open browser")
	}

	return err
}

func showVersion() {
	fmt.Println("mymine version", version)
}

type Options struct {
	Open    []int  `short:"o" long:"open"    description:"Open specified ticket on a web browser"`
	Version []bool `short:"v" long:"version" description:"Show version"`
}

func main() {
	var opts Options
	parser := flags.NewParser(&opts, flags.Default)

	parser.Name = "mymine"
	parser.Usage = "[OPTIONS]"

	_, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}

	if opts.Version != nil {
		showVersion()
		os.Exit(1)
	}

	redmineURL := lookupEnv("REDMINE_URL")
	if redmineURL == "" {
		fmt.Println("REDMINE_URL is not specified.")
		os.Exit(1)
	}

	if opts.Open != nil {
		url := redmineURL + "issues/" + strconv.Itoa(opts.Open[0])
		openURLByBrowser(url)
		os.Exit(1)
	}

	redmineApiKey := lookupEnv("REDMINE_API_KEY")
	if redmineApiKey == "" {
		fmt.Println("REDMINE_API_KEY is not specified.")
		os.Exit(1)
	}

	request := redmineURL + "issues.json?key=" + redmineApiKey + "&status_id=open&assigned_to_id=me&limit=100"
	fmt.Println("request =", request)
	fmt.Println("fetching information...")
	var buf map[string]interface{}
	rest.Get(&buf, request, nil)

	issues := buf["issues"].([]interface{})
	for _, v := range issues {
		issue := v.(map[string]interface{})

		id := int(issue["id"].(float64))
		status := issue["status"].(map[string]interface{})
		fmt.Printf("[#%d] %11s %s\n", id, status["name"], issue["subject"])
	}
}
