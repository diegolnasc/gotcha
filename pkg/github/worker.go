// Copyright 2021 Diego Lima. All rights reserved.

// Use of this source code is governed by a Apache license.
// license that can be found in the LICENSE file.

package github

import (
	"strconv"
	"strings"

	ghwebhooks "github.com/go-playground/webhooks/v6/github"
	v41 "github.com/google/go-github/v41/github"
)

// processPullRequestEvent start a pull request event process.
func (s *PullRequestService) processPullRequestEvent(p *ghwebhooks.PullRequestPayload) {
	owner, _ := getOwner(&s.w.Config)
	s.processPullRequest(owner, p)
}

// processIssueCommentEvent start an issue comment event process.
func (s *IssueService) processIssueCommentEvent(p *ghwebhooks.IssueCommentPayload) {
	if isUserAuthorized(&s.w.Config, &p.Sender.Login, &p.Repository.Name) && p.Action == "created" {
		var pullRequest *v41.PullRequest
		owner, _ := getOwner(&s.w.Config)
		if len(p.Issue.PullRequest.HTMLURL) > 0 {
			if pullRequestNumber, err := strconv.Atoi(strings.Split(p.Issue.PullRequest.HTMLURL, "/")[6]); err == nil {
				pullRequest, err = s.w.GetPullRequest(*owner, p.Repository.Name, pullRequestNumber)
				if err != nil {
					return
				}
			}
		}
		s.processIssueComment(owner, pullRequest, p)
	}
}

// processCheckRunEvent start an check run event process.
func (s *CheckService) processCheckRunEvent(p *ghwebhooks.CheckRunPayload) {
	if p.CheckRun.App.ID == int64(s.w.Config.Github.AppID) {
		var pullRequest *v41.PullRequest
		owner, _ := getOwner(&s.w.Config)
		if len(p.CheckRun.PullRequests) > 0 {
			pullRequest, _ = s.w.GetPullRequest(*owner, p.Repository.Name, int(p.CheckRun.PullRequests[0].Number))
		} else {
			chekRun, err := s.w.GetCheckRun(*owner, p.Repository.Name, p.CheckRun.ID)
			if err != nil {
				return
			}
			if len(*chekRun.ExternalID) > 0 {
				pullRequestID, _ := strconv.Atoi(*chekRun.ExternalID)
				pullRequest, _ = s.w.GetPullRequest(*owner, p.Repository.Name, pullRequestID)
			}
		}
		s.processCheckRun(owner, pullRequest, p)
	}
}
