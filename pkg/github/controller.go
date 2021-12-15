package github

import (
	"context"
	"log"

	v41 "github.com/google/go-github/v41/github"
)

func (w *Worker) PullRequestCreateReview(owner string, repo string, number int, review v41.PullRequestReviewRequest) (*v41.PullRequestReview, error) {
	resp, _, err := w.Client.PullRequests.CreateReview(context.TODO(), owner, repo, number, &review)
	if err != nil {
		log.Printf("error creating review: %v\n", err)
	}
	return resp, err
}

func (w *Worker) PullRequestCreateCommentReview(owner string, repo string, number int, comment v41.PullRequestComment) (*v41.PullRequestComment, error) {
	resp, _, err := w.Client.PullRequests.CreateComment(context.TODO(), owner, repo, number, &comment)
	if err != nil {
		log.Printf("error creating pull request review comment: %v\n", err)
	}
	return resp, err
}

func (w *Worker) PullRequestCreateComment(owner string, repo string, number int, comment v41.IssueComment) (*v41.IssueComment, error) {
	resp, _, err := w.Client.Issues.CreateComment(context.TODO(), owner, repo, number, &comment)
	if err != nil {
		log.Printf("error creating pull request comment: %v\n", err)
	}
	return resp, err
}

func (w *Worker) AddLabels(owner string, repo string, number int, labels []string) ([]*v41.Label, error) {
	resp, _, err := w.Client.Issues.AddLabelsToIssue(context.TODO(), owner, repo, number, labels)
	if err != nil {
		log.Printf("error adding labels: %v\n", err)
	}
	return resp, err
}

func (w *Worker) CreateCheckRun(owner string, repo string, checkrun v41.CreateCheckRunOptions) (*v41.CheckRun, error) {
	resp, _, err := w.Client.Checks.CreateCheckRun(context.TODO(), owner, repo, checkrun)
	if err != nil {
		log.Printf("error creating checkrun: %v\n", err)
	}
	return resp, err
}

func (w *Worker) UpdateCheckRun(owner string, repo string, checkRunId int64, checkRun v41.UpdateCheckRunOptions) (*v41.CheckRun, error) {
	resp, _, err := w.Client.Checks.UpdateCheckRun(context.TODO(), owner, repo, checkRunId, checkRun)
	if err != nil {
		log.Printf("error updating checkrun: %v\n", err)
	}
	return resp, err
}

func (w *Worker) GetPullRequest(owner string, repo string, number int) (*v41.PullRequest, error) {
	resp, _, err := w.Client.PullRequests.Get(context.TODO(), owner, repo, number)
	if err != nil {
		log.Printf("error getting pull request: %v\n", err)
	}
	return resp, err
}
