package github

import (
	"strconv"

	ghwebhooks "github.com/go-playground/webhooks/v6/github"
	v41 "github.com/google/go-github/v41/github"
)

type PullRequestService service

func (s *PullRequestService) processPullRequest(owner *string, p *ghwebhooks.PullRequestPayload) {
	s.createPullRequestOverview(owner, p)
	if p.Action == "opened" || p.Action == "reopened" || p.Action == "edited" {
		pullRequestNumber := strconv.Itoa(int(p.PullRequest.Number))
		s.w.CreateCheckRun(*owner, p.Repository.Name, v41.CreateCheckRunOptions{
			Name:       "Laboratory test",
			Status:     v41.String("in_progress"),
			HeadSHA:    p.PullRequest.Head.Sha,
			DetailsURL: &p.PullRequest.HTMLURL,
			ExternalID: &pullRequestNumber,
		})
	}
}

func (s *PullRequestService) createPullRequestOverview(owner *string, p *ghwebhooks.PullRequestPayload) {
	if (p.Action == "opened" || p.Action == "reopened") && s.w.Config.Layout.PullRequest.EnableOverview {
		s.w.CreatePulllRequestOverviewComment(owner, p.Repository.Name, int(p.PullRequest.Number))
	}
}
