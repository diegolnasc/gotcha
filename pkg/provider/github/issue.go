// Copyright 2021 Diego Lima. All rights reserved.

// Use of this source code is governed by a Apache license.
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	ghwebhooks "github.com/go-playground/webhooks/v6/github"
	v41 "github.com/google/go-github/v41/github"
)

// IssueService handles communication with the issue event.
// Issue Comment also refer to a pull request comment.
type IssueService service

// processIssueComment process the issue comment event payload.
func (s *IssueService) processIssueComment(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.IssueCommentPayload) {
	s.reRunPullRequestOverview(owner, pullRequest, p)
	s.approve(owner, pullRequest, p)
	s.merge(owner, pullRequest, p)
	s.mergeAndDelete(owner, pullRequest, p)
	s.reRunLaboratoryTest(owner, pullRequest, p)
}

// reRunPullRequestOverview re-run the pull request report (diff overview).
func (s *IssueService) reRunPullRequestOverview(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.IssueCommentPayload) {
	if pullRequest != nil {
		if strings.EqualFold(p.Comment.Body, s.w.Config.Layout.PullRequest.OverViewCommand) {
			s.w.CreatePulllRequestOverviewComment(owner, p.Repository.Name, *pullRequest.Number)
		}
	}
}

// approve approves a pull request.
func (s *IssueService) approve(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.IssueCommentPayload) {
	if pullRequest != nil {
		if strings.EqualFold(p.Comment.Body, s.w.Config.Layout.PullRequest.ApproveCommand) {
			message := fmt.Sprintf("[%s] Looks Good To Me!", p.Issue.User.Login)
			_, err := s.w.PullRequestCreateReview(*owner, p.Repository.Name, *pullRequest.Number, v41.PullRequestReviewRequest{
				Body:  &message,
				Event: v41.String("APPROVE"),
			})
			if err != nil {
				log.Printf("error creview: %v\n", err)
			}
		}
	}
}

// merge a pull request.
func (s *IssueService) merge(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.IssueCommentPayload) {
	if pullRequest != nil {
		if strings.EqualFold(p.Comment.Body, s.w.Config.Layout.PullRequest.MergeCommand) {
			s.w.MergePullRequest(*owner, p.Repository.Name, *pullRequest.Number, "", v41.PullRequestOptions{})
		}
	}
}

// mergeAndDelete merge and delete the reference of the pull request.
func (s *IssueService) mergeAndDelete(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.IssueCommentPayload) {
	if pullRequest != nil {
		if strings.EqualFold(p.Comment.Body, s.w.Config.Layout.PullRequest.MergeAndDeleteCommand) {
			if _, err := s.w.MergePullRequest(*owner, p.Repository.Name, *pullRequest.Number, "", v41.PullRequestOptions{}); err == nil {
				if *pullRequest.Base.Repo.Name == *pullRequest.Head.Repo.Name {
					s.w.DeleteRef(*owner, p.Repository.Name, fmt.Sprintf("heads/%s", *pullRequest.Head.Ref))
				}
			}
		}
	}
}

// reRunLaboratoryTest re-run the pull request test suite.
func (s *IssueService) reRunLaboratoryTest(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.IssueCommentPayload) {
	if pullRequest != nil {
		if strings.EqualFold(p.Comment.Body, s.w.Config.Layout.PullRequest.RunTestSuiteCommand) {
			pullRequestNumber := strconv.Itoa(int(*pullRequest.Number))
			s.w.CreateCheckRun(*owner, p.Repository.Name, v41.CreateCheckRunOptions{
				Name:       "Laboratory test",
				Status:     v41.String("in_progress"),
				HeadSHA:    *pullRequest.Head.SHA,
				DetailsURL: pullRequest.HTMLURL,
				ExternalID: &pullRequestNumber,
			})
		}
	}
}
