package github

import (
	"bytes"
	"fmt"
	"regexp"

	ghwebhooks "github.com/go-playground/webhooks/v6/github"
	v41 "github.com/google/go-github/v41/github"
)

type CheckService service

type checkRunResult struct {
	title  string
	passed bool
	body   string
}

func (s *CheckService) processCheckRun(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.CheckRunPayload) {
	if p.Action == "created" {
		result, overralResults := s.runTestSuite(owner, pullRequest, p)
		s.updateCheckRunStatus(owner, pullRequest, p, overralResults)
		s.printResults(owner, pullRequest, p, result, overralResults)
	}
}

func (s *CheckService) runTestSuite(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.CheckRunPayload) ([]*checkRunResult, bool) {
	var results []*checkRunResult
	overralResults := true
	if len(s.w.Config.Layout.PullRequest.TestSuite.NamePattern) > 0 {
		results = append(results, s.isNamePatternValid(owner, pullRequest, p))
	}
	if s.w.Config.Layout.PullRequest.TestSuite.Reviewers {
		results = append(results, s.hasReviewers(owner, pullRequest, p))
	}
	if s.w.Config.Layout.PullRequest.TestSuite.Assignees {
		results = append(results, s.hasAssignees(owner, pullRequest, p))
	}
	if s.w.Config.Layout.PullRequest.TestSuite.Labels {
		results = append(results, s.hasLabels(owner, pullRequest, p))
	}
	for _, r := range results {
		if !r.passed {
			overralResults = false
			break
		}
	}
	return results, overralResults
}

func (s *CheckService) isNamePatternValid(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.CheckRunPayload) *checkRunResult {
	match, _ := regexp.MatchString(s.w.Config.Layout.PullRequest.TestSuite.NamePattern, *pullRequest.Title)
	var body string
	if !match {
		body = fmt.Sprintf("The pull request format should be: [%s]", s.w.Config.Layout.PullRequest.TestSuite.NamePattern)
	}
	return &checkRunResult{
		title:  "Pull request pattern",
		passed: match,
		body:   body,
	}
}

func (s *CheckService) hasReviewers(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.CheckRunPayload) *checkRunResult {
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

func (s *CheckService) hasAssignees(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.CheckRunPayload) *checkRunResult {
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

func (s *CheckService) hasLabels(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.CheckRunPayload) *checkRunResult {
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

func (s *CheckService) updateCheckRunStatus(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.CheckRunPayload, overralResults bool) {
	var conclusion string
	if overralResults {
		conclusion = "success"
	} else {
		conclusion = "failure"
	}
	s.w.UpdateCheckRun(*owner, p.Repository.Name, p.CheckRun.ID, v41.UpdateCheckRunOptions{
		Status:     v41.String("completed"),
		Conclusion: &conclusion,
	})
}

func (s *CheckService) printResults(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.CheckRunPayload, results []*checkRunResult, overralResults bool) {
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
	s.w.IssueCreateComment(*owner, p.Repository.Name, *pullRequest.Number, v41.IssueComment{
		Body: v41.String(body.String()),
	})
}
