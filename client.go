package main

import (
	"net/http"
)

type client struct {
	endpoint string
	apikey   string
	*http.Client
}

func newClient(endpoint, apikey string) *client {
	return &client{endpoint, apikey, http.DefaultClient}
}

type errorsResult struct {
	Errors []string `json:"errors"`
}
