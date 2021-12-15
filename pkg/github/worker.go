package github

import (
	"strconv"
	"strings"

	ghwebhooks "github.com/go-playground/webhooks/v6/github"
	v41 "github.com/google/go-github/v41/github"
)

func (w *Worker) processPullRequestEvent(p *ghwebhooks.PullRequestPayload) {
	owner, _ := getOwner(&w.Config)
	w.processPullRequest(owner, p)
}

func (w *Worker) processIssueCommentEvent(p *ghwebhooks.IssueCommentPayload) {
	if isUserAuthorized(&w.Config, &p.Sender.Login, &p.Repository.Name) {
		var pullRequest *v41.PullRequest
		owner, _ := getOwner(&w.Config)
		if len(p.Issue.PullRequest.HTMLURL) > 0 {
			pullRequestId, err := strconv.Atoi(strings.Split(p.Issue.PullRequest.HTMLURL, "/")[6])
			if err == nil {
				pullRequest, err = w.GetPullRequest(*owner, p.Repository.Name, pullRequestId)
				if err != nil {
					return
				}
			}
		}
		w.processIssueComment(owner, pullRequest, p)
	}
}

func (w *Worker) processCheckRunEvent(p *ghwebhooks.CheckRunPayload) {
	owner, _ := getOwner(&w.Config)
	pullRequest, err := w.GetPullRequest(*owner, p.Repository.Name, int(p.CheckRun.PullRequests[0].Number))
	if err != nil {
		return
	}
	if p.CheckRun.App.ID == int64(w.Config.Github.AppID) {
		w.processCheckRun(owner, pullRequest, p)
	}
}
