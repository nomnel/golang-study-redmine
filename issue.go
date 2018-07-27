package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type issue struct {
	ID          int    `json:"id"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
	ProjectID   int    `json:"project_id"`
}

type issueFilter struct {
	ProjectID    string
	SubprojectID string
	TrackerID    string
	StatusID     string
	AssignedToID string
	UpdatedOn    string
}

type issueRequest struct {
	Issue issue `json:"issue"`
}

type issuesResult struct {
	Issues     []issue `json:"issues"`
	TotalCount uint    `json:"total_count"`
	Offset     uint    `json:"offset"`
	Limit      uint    `json:"limit"`
}

// issue が値渡しなのはなぜ？
func (c *client) createIssue(i issue) (*issue, error) {
	ir := issueRequest{
		Issue: i,
	}
	s, err := json.Marshal(ir)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.endpoint+"/issues.json?key="+c.apikey, strings.NewReader(string(s)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r issueRequest
	if res.StatusCode != 201 {
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
	return &r.Issue, nil
}

func (c *client) issues() ([]issue, error) {
	issues, err := getIssues(c, "/issues.json?key="+c.apikey)

	if err != nil {
		return nil, err
	}

	return issues, nil
}

func (c *client) issuesByFilter(f *issueFilter) ([]issue, error) {
	issues, err := getIssues(c, "/issues.json?key="+c.apikey+c.getPaginationClause()+getIssueFilterClause(f))

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

// issueFilter.toClause じゃだめなのか？
func getIssueFilterClause(filter *issueFilter) string {
	if filter == nil {
		return ""
	}

	clause := ""
	if filter.ProjectID != "" {
		clause = clause + fmt.Sprintf("&project_id=%v", filter.ProjectID)
	}
	if filter.SubprojectID != "" {
		clause = clause + fmt.Sprintf("&subproject_id=%v", filter.SubprojectID)
	}
	if filter.TrackerID != "" {
		clause = clause + fmt.Sprintf("&tracker_id=%v", filter.TrackerID)
	}
	if filter.StatusID != "" {
		clause = clause + fmt.Sprintf("&status_id=%v", filter.StatusID)
	}
	if filter.AssignedToID != "" {
		clause = clause + fmt.Sprintf("&assigned_to_id=%v", filter.AssignedToID)
	}
	if filter.UpdatedOn != "" {
		clause = clause + fmt.Sprintf("&updated_on=%v", filter.UpdatedOn)
	}

	return clause
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
