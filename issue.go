package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type issue struct {
	Id      int    `json:"id"`
	Subject string `json:"subject"`
}

type issuesResult struct {
	Issues     []issue `json:"issues"`
	TotalCount uint    `json:"total_count"`
	Offset     uint    `json:"offset"`
	Limit      uint    `json:"limit"`
}

func (c *client) issues() ([]issue, error) {
	issues, err := getIssues(c, "/issues.json?key="+c.apikey)

	if err != nil {
		return nil, err
	}

	return issues, nil
}

func getIssue(c *client, url string, offset int) (*issuesResult, error) {
	res, err := c.Get(c.endpoint + url + "&offset=" + strconv.Itoa(offset))

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r issuesResult
	if res.StatusCode != 200 {
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	} else {
		err = decoder.Decode(&r)
	}
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func getIssues(c *client, url string) ([]issue, error) {
	completed := false
	var issues []issue

	for completed == false {
		r, err := getIssue(c, url, len(issues))

		if err != nil {
			return nil, err
		}

		if r.TotalCount == uint(len(issues)) {
			completed = true
		}

		issues = append(issues, r.Issues...)
	}

	return issues, nil
}
