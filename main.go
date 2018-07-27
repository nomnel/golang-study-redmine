package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type config struct {
	Endpoint string `json:"endpoint"`
	Apikey   string `json:"apikey"`
	Project  int    `json:"project"`
}

var conf config

func main() {
	conf = getConfig()

	flag.Parse()
	switch flag.Arg(0) {
	case "i", "issue":
		switch flag.Arg(1) {
		case "a", "add":
			if flag.NArg() == 4 {
				addIssue(flag.Arg(2), flag.Arg(3))
			} else {
				usage()
			}
			break
		case "d", "delete":
			if flag.NArg() == 3 {
				id, err := strconv.Atoi(flag.Arg(2))
				if err != nil {
					fatal("Invalid issue id: %s\n", err)
				}
				deleteIssue(id)
			} else {
				usage()
			}
			break
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

func addIssue(subject, description string) {
	c := newClient(conf.Endpoint, conf.Apikey)
	i := issue{
		ProjectID:   conf.Project,
		Subject:     subject,
		Description: description,
	}
	_, err := c.createIssue(i)
	if err != nil {
		fatal("Failed to create issue: %s\n", err)
	}
}

func deleteIssue(id int) {
	c := newClient(conf.Endpoint, conf.Apikey)
	err := c.deleteIssue(id)
	if err != nil {
		fatal("Failed to delete issue: %s\n", err)
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
