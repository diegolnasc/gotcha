// Copyright 2021 Diego Lima. All rights reserved.

// Use of this source code is governed by a Apache license.
// license that can be found in the LICENSE file.

package github

import (
	"strconv"

	ghwebhooks "github.com/go-playground/webhooks/v6/github"
	v41 "github.com/google/go-github/v41/github"
)

// PullRequestService handles communication with the pull request event.
type PullRequestService service

// processPullRequest process the pull request event payload.
func (s *PullRequestService) processPullRequest(owner *string, p *ghwebhooks.PullRequestPayload) {
	s.createPullRequestOverview(owner, p)
	if p.Action == "opened" || p.Action == "reopened" || p.Action == "edited" || p.Action == "synchronize" {
		if p.Action == "edited" && p.Changes.Title == nil {
		} else {
			pullRequestNumber := strconv.Itoa(int(p.PullRequest.Number))
			if _, err := s.w.CreateCheckRun(*owner, p.Repository.Name, v41.CreateCheckRunOptions{
				Name:       "Laboratory test",
				Status:     v41.String("in_progress"),
				HeadSHA:    p.PullRequest.Head.Sha,
				DetailsURL: &p.PullRequest.HTMLURL,
				ExternalID: &pullRequestNumber,
			}); err != nil {
				//erro
			}
		}
	}
}

// createPullRequestOverview create the pull request report (diff overview).
func (s *PullRequestService) createPullRequestOverview(owner *string, p *ghwebhooks.PullRequestPayload) {
	if (p.Action == "opened" || p.Action == "reopened") && s.w.Config.Layout.PullRequest.EnableOverview {
		s.w.CreatePulllRequestOverviewComment(owner, p.Repository.Name, int(p.PullRequest.Number))
	}
}
