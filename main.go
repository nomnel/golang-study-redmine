package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type config struct {
	Endpoint string `json:"endpoint"`
	Apikey   string `json:"apikey"`
}

func main() {
	c := getConfig()
	client := newClient(c.Endpoint, c.Apikey)
	issues, err := client.issues()
	if err != nil {
		fatal("Failed to list issues: %s\n", err)
	}
	for _, i := range issues {
		fmt.Printf("%4d: %s\n", i.Id, i.Subject)
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
