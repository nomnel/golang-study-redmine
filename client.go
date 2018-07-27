package main

import (
	"fmt"
	"net/http"
)

type client struct {
	endpoint string
	apikey   string
	*http.Client
	Limit  int
	Offset int
}

var defaultLimit = -1  // "-1" means "No setting"
var defaultOffset = -1 //"-1" means "No setting"

func newClient(endpoint, apikey string) *client {
	return &client{endpoint, apikey, http.DefaultClient, defaultLimit, defaultOffset}
}

type errorsResult struct {
	Errors []string `json:"errors"`
}

func (c *client) getPaginationClause() string {
	clause := ""
	if c.Limit > -1 {
		clause = clause + fmt.Sprintf("&limit=%v", c.Limit)
	}
	if c.Offset > -1 {
		clause = clause + fmt.Sprintf("&offset=%v", c.Offset)
	}
	return clause
}
