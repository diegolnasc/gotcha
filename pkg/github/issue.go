package github

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	ghwebhooks "github.com/go-playground/webhooks/v6/github"
	v41 "github.com/google/go-github/v41/github"
)

func (w *Worker) processIssueComment(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.IssueCommentPayload) {
	w.approve(owner, pullRequest, p)
	w.merge(owner, pullRequest, p)
	w.mergeAndDelete(owner, pullRequest, p)
	w.reRunLaboratoryTest(owner, pullRequest, p)
}

func (w *Worker) approve(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.IssueCommentPayload) {
	if pullRequest != nil {
		if strings.EqualFold(p.Comment.Body, w.Config.Layout.PullRequest.ApproveCommand) {
			message := fmt.Sprintf("[%s] Looks Good To Me!", p.Issue.User.Login)
			_, err := w.PullRequestCreateReview(*owner, p.Repository.Name, *pullRequest.Number, v41.PullRequestReviewRequest{
				Body:  &message,
				Event: v41.String("APPROVE"),
			})
			if err != nil {
				log.Printf("error creview: %v\n", err)
			}
		}
	}
}

func (w *Worker) merge(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.IssueCommentPayload) {
	if pullRequest != nil {
		if strings.EqualFold(p.Comment.Body, w.Config.Layout.PullRequest.MergeCommand) {
			w.MergePullRequest(*owner, p.Repository.Name, *pullRequest.Number, "", v41.PullRequestOptions{})
		}
	}
}

func (w *Worker) mergeAndDelete(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.IssueCommentPayload) {
	if pullRequest != nil {
		if strings.EqualFold(p.Comment.Body, w.Config.Layout.PullRequest.MergeAndDeleteCommand) {
			if _, err := w.MergePullRequest(*owner, p.Repository.Name, *pullRequest.Number, "", v41.PullRequestOptions{}); err == nil {
				if *pullRequest.Base.Repo.Name == *pullRequest.Head.Repo.Name {
					w.DeleteRef(*owner, p.Repository.Name, fmt.Sprintf("heads/%s", *pullRequest.Head.Ref))
				}
			}
		}
	}
}

func (w *Worker) reRunLaboratoryTest(owner *string, pullRequest *v41.PullRequest, p *ghwebhooks.IssueCommentPayload) {
	if pullRequest != nil {
		if strings.EqualFold(p.Comment.Body, w.Config.Layout.PullRequest.RunTestSuiteCommand) {
			pullRequestNumber := strconv.Itoa(int(*pullRequest.Number))
			w.CreateCheckRun(*owner, p.Repository.Name, v41.CreateCheckRunOptions{
				Name:       "Laboratory test",
				Status:     v41.String("in_progress"),
				HeadSHA:    *pullRequest.Head.SHA,
				DetailsURL: pullRequest.HTMLURL,
				ExternalID: &pullRequestNumber,
			})
		}
	}
}
