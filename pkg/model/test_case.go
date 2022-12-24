// Copyright 2021 Diego Lima. All rights reserved.

// Use of this source code is governed by a Apache license.
// license that can be found in the LICENSE file.

package model

import (
	"bytes"
	"fmt"
	"regexp"
)

// Result represents the result of a test case.
type Result struct {
	Title  string
	Passed bool
	Body   string
}

// IsNamePatternValid check if the pull request name is valid.
func (cfg *Layout) IsNamePatternValid(name string) *Result {
	fmt.Println(cfg)
	if len(cfg.PullRequest.TestSuite.NamePattern) > 0 {
		match, _ := regexp.MatchString(cfg.PullRequest.TestSuite.NamePattern, name)
		var body string
		if !match {
			body = fmt.Sprintf("The pull request format should be: [%s]", cfg.PullRequest.TestSuite.NamePattern)
		}
		return &Result{
			Title:  "Pull request pattern",
			Passed: match,
			Body:   body,
		}
	}
	return &Result{}
}

// HasReviewers check if the pull request has reviewers.
func (cfg *Layout) HasReviewers(hasReviewers bool) *Result {
	if cfg.PullRequest.TestSuite.Reviewers {
		if !hasReviewers {
			return &Result{
				Title:  "Reviewers",
				Passed: false,
				Body:   "Must have at least one reviewer",
			}
		}
		return &Result{
			Title:  "Reviewers",
			Passed: true,
		}
	}
	return &Result{}
}

// HasAssignees check if the pull request has assignees.
func (cfg *Layout) HasAssignees(hasAssignees bool) *Result {
	if cfg.PullRequest.TestSuite.Assignees {
		if !hasAssignees {
			return &Result{
				Title:  "Assignees",
				Passed: false,
				Body:   "Must have at least one assignee",
			}
		}
		return &Result{
			Title:  "Assignees",
			Passed: true,
		}
	}
	return &Result{}
}

// HasLabels check if the pull request has labels.
func (cfg *Layout) HasLabels(hasLabels bool) *Result {
	if cfg.PullRequest.TestSuite.Assignees {
		if !hasLabels {
			return &Result{
				Title:  "Labels",
				Passed: false,
				Body:   "Must have at least one label",
			}
		}
		return &Result{
			Title:  "Labels",
			Passed: true,
		}
	}
	return &Result{}
}

func (cfg *Layout) GetOverralResults(results []*Result) string {
	var title string
	var body bytes.Buffer
	overralResults := cfg.AllResultsPassed(results)
	if overralResults {
		title = `<h3 align="center">Laboratory test results [:heavy_check_mark:]</h3>

---
`
	} else {
		title = `<h3 align="center">Laboratory test results [:x:]</h3>

---
`
	}
	body.WriteString(title)
	for _, r := range results {
		if r.Passed {
			body.WriteString(fmt.Sprintf("[:white_check_mark:] **%s** <br/>", r.Title))
		} else {
			body.WriteString(fmt.Sprintf("[:heavy_exclamation_mark:] **%s** &#8594; %s <br/>", r.Title, r.Body))
		}
	}
	return body.String()
}

// AllResultsPassed check if all the results passed.
func (cfg *Layout) AllResultsPassed(results []*Result) bool {
	overralResults := true
	for _, r := range results {
		if !r.Passed {
			overralResults = false
			break
		}
	}
	return overralResults
}
