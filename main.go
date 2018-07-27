package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type config struct {
	Endpoint string `json:"endpoint"`
	Apikey   string `json:"apikey"`
}

var conf config

func main() {
	conf = getConfig()

	flag.Parse()
	switch flag.Arg(0) {
	case "i", "issue":
		switch flag.Arg(1) {
		case "l", "list":
			listIssues(nil)
			break
		case "m", "mine":
			filter := issueFilter{
				AssignedToID: "me",
			}
			listIssues(&filter)
			break
		default:
			usage()
		}
	default:
		usage()
	}
}

func getConfig() config {
	filename := "settings.json"
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fatal("Failed to read config file: %s\n", err)
	}

	var c config
	err = json.Unmarshal(b, &c)
	if err != nil {
		fatal("Failed to unmarshal file: %s\n", err)
	}
	return c
}

func fatal(format string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, format, err)
	} else {
		fmt.Fprint(os.Stderr, format)
	}
	os.Exit(1)
}

func listIssues(filter *issueFilter) {
	c := newClient(conf.Endpoint, conf.Apikey)
	issues, err := c.issuesByFilter(filter)
	if err != nil {
		fatal("Failed to list issues: %s\n", err)
	}
	for _, i := range issues {
		fmt.Printf("%4d: %s\n", i.ID, i.Subject)
	}
}

func usage() {
	fmt.Println(`./golang-study-redmine <command>

Issue Commands:
  list     l listing issues.
             $ ./golang-study-redmine i l
  mine     m listing your issues.
             $ ./golang-study-redmine i m`)
	os.Exit(1)
}
