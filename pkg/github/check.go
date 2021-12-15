package github

import (
	"bytes"
	"fmt"
	"regexp"

	ghwebhooks "github.com/go-playground/webhooks/v6/github"
	v41 "github.com/google/go-github/v41/github"
)

type checkRunResult struct {
	title  string
	passed bool
	body   string
}

func (w *Worker) processCheckRun(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.CheckRunPayload) {
	if p.Action == "created" {
		result, overralResults := w.runTestSuite(owner, pullRequest, p)
		w.updateCheckRunStatus(owner, pullRequest, p, overralResults)
		w.printResults(owner, pullRequest, p, result, overralResults)
	}
}

func (w *Worker) runTestSuite(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.CheckRunPayload) ([]*checkRunResult, bool) {
	var results []*checkRunResult
	overralResults := true
	if len(w.Config.Layout.PullRequest.TestSuite.NamePattern) > 0 {
		results = append(results, w.isNamePatternValid(owner, pullRequest, p))
	}
	if w.Config.Layout.PullRequest.TestSuite.Reviewers {
		results = append(results, w.hasReviewers(owner, pullRequest, p))
	}
	if w.Config.Layout.PullRequest.TestSuite.Assignees {
		results = append(results, w.hasAssignees(owner, pullRequest, p))
	}
	if w.Config.Layout.PullRequest.TestSuite.Labels {
		results = append(results, w.hasLabels(owner, pullRequest, p))
	}
	for _, r := range results {
		if !r.passed {
			overralResults = false
			break
		}
	}
	return results, overralResults
}

func (w *Worker) isNamePatternValid(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.CheckRunPayload) *checkRunResult {
	match, _ := regexp.MatchString(w.Config.Layout.PullRequest.TestSuite.NamePattern, *pullRequest.Title)
	var body string
	if !match {
		body = fmt.Sprintf("The pull request format should be: [%s]", w.Config.Layout.PullRequest.TestSuite.NamePattern)
	}
	return &checkRunResult{
		title:  "Pull request pattern",
		passed: match,
		body:   body,
	}
}

func (w *Worker) hasReviewers(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.CheckRunPayload) *checkRunResult {
	if len(pullRequest.RequestedReviewers) == 0 {
		return &checkRunResult{
			title:  "Reviewers",
			passed: false,
			body:   "Must have at least one reviewer",
		}
	}
	return &checkRunResult{
		title:  "Reviewers",
		passed: true,
	}
}

func (w *Worker) hasAssignees(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.CheckRunPayload) *checkRunResult {
	if len(pullRequest.Assignees) == 0 {
		return &checkRunResult{
			title:  "Assignees",
			passed: false,
			body:   "Must have at least one assignee",
		}
	}
	return &checkRunResult{
		title:  "Assignees",
		passed: true,
	}
}

func (w *Worker) hasLabels(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.CheckRunPayload) *checkRunResult {
	if len(pullRequest.Labels) == 0 {
		return &checkRunResult{
			title:  "Labels",
			passed: false,
			body:   "Must have at least one label",
		}
	}
	return &checkRunResult{
		title:  "Labels",
		passed: true,
	}
}

func (w *Worker) updateCheckRunStatus(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.CheckRunPayload, overralResults bool) {
	var conclusion string
	if overralResults {
		conclusion = "success"
	} else {
		conclusion = "failure"
	}
	w.UpdateCheckRun(*owner, p.Repository.Name, p.CheckRun.ID, v41.UpdateCheckRunOptions{
		Status:     v41.String("completed"),
		Conclusion: &conclusion,
	})
}

func (w *Worker) printResults(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.CheckRunPayload, results []*checkRunResult, overralResults bool) {
	var title string
	var body bytes.Buffer
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
		if r.passed {
			body.WriteString(fmt.Sprintf("[:white_check_mark:] **%s** <br/>", r.title))
		} else {
			body.WriteString(fmt.Sprintf("[:heavy_exclamation_mark:] **%s** &#8594; %s <br/>", r.title, r.body))
		}
	}
	w.PullRequestCreateComment(*owner, p.Repository.Name, *pullRequest.Number, v41.IssueComment{
		Body: v41.String(body.String()),
	})
}
