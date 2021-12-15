package github

import (
	ghwebhooks "github.com/go-playground/webhooks/v6/github"
	v41 "github.com/google/go-github/v41/github"
)

func (w *Worker) processPullRequest(owner *string, p *ghwebhooks.PullRequestPayload) {
	if p.Action == "opened" || p.Action == "reopened" || p.Action == "edited" {
		w.CreateCheckRun(*owner, p.Repository.Name, v41.CreateCheckRunOptions{
			Name:    "Laboratory test",
			Status:  v41.String("in_progress"),
			HeadSHA: p.PullRequest.Head.Sha,
		})
	}
}
