// Copyright 2021 Diego Lima. All rights reserved.

// Use of this source code is governed by a Apache license.
// license that can be found in the LICENSE file.

package github

import (
	"log"

	"github.com/diegolnasc/gotcha/pkg/model"
	ghwebhooks "github.com/go-playground/webhooks/v6/github"
	v41 "github.com/google/go-github/v41/github"
)

// CheckService handles communication with the checkrun/checksuite event.
type CheckService service

// processCheckRun process the checkrun event payload.
func (s *CheckService) processCheckRun(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.CheckRunPayload) {
	if p.Action == "created" {
		result, overralResults := s.runTestSuite(owner, pullRequest, p)
		s.updateCheckRunStatus(owner, pullRequest, p, overralResults)
		s.printResults(owner, pullRequest, p, result, overralResults)
	}
}

// runTestSuite run the test suite.
func (s *CheckService) runTestSuite(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.CheckRunPayload) ([]*model.Result, bool) {
	var results []*model.Result
	results = append(results, s.w.Config.Layout.IsNamePatternValid(*pullRequest.Title))
	results = append(results, s.w.Config.Layout.HasReviewers(len(pullRequest.RequestedReviewers) > 0))
	results = append(results, s.w.Config.Layout.HasAssignees(len(pullRequest.Assignees) > 0))
	results = append(results, s.w.Config.Layout.HasLabels(len(pullRequest.Labels) > 0))
	return results, s.w.Config.Layout.AllResultsPassed(results)
}

// updateCheckRunStatus update the check run status in github.
func (s *CheckService) updateCheckRunStatus(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.CheckRunPayload, overralResults bool) {
	var conclusion string
	if overralResults {
		conclusion = "success"
	} else {
		conclusion = "failure"
	}
	if _, err := s.w.UpdateCheckRun(*owner, p.Repository.Name, p.CheckRun.ID, v41.UpdateCheckRunOptions{
		Status:     v41.String("completed"),
		Conclusion: &conclusion,
	}); err != nil {
		log.Printf("Unable to complete the provider call %s.", err)
	}
}

// printResults comments the result of the test suite on the pull request.
func (s *CheckService) printResults(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.CheckRunPayload, results []*model.Result, overralResults bool) {
	if _, err := s.w.IssueCreateComment(*owner, p.Repository.Name, *pullRequest.Number, v41.IssueComment{
		Body: v41.String(s.w.Config.Layout.GetOverralResults(results)),
	}); err != nil {
		log.Printf("Unable to complete the provider call %s.", err)
	}
}
