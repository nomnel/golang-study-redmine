package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		panic("Panic!")
	}

	var c config
	err = json.Unmarshal(b, &c)
	if err != nil {
		panic("Panic!")
	}
	return c
}
