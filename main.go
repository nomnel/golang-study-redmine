package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type config struct {
	Endpoint string `json:"endpoint"`
}

func main() {
	c := getConfig()
	fmt.Printf("%s", c.Endpoint)
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
