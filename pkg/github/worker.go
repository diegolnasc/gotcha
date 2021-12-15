package github

import (
	ghwebhooks "github.com/go-playground/webhooks/v6/github"
)

func (w *Worker) processPullRequestEvent(p *ghwebhooks.PullRequestPayload) {
	// w.processPullRequest(p)
}

func (w *Worker) processIssueCommentEvent(p *ghwebhooks.IssueCommentPayload) {
	// w.processIssueComment(p)
}

func (w *Worker) processCheckRunEvent(p *ghwebhooks.CheckRunPayload) {
	// w.processCheckRun(p)
}
