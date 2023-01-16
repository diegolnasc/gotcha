// Copyright 2021 Diego Lima. All rights reserved.

// Use of this source code is governed by a Apache license.
// license that can be found in the LICENSE file.

package github

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	v41 "github.com/google/go-github/v41/github"
)

// GetPullRequest get a pull request.
func (w *Worker) GetPullRequest(owner string, repo string, number int) (*v41.PullRequest, error) {
	resp, _, err := w.Client.PullRequests.Get(context.TODO(), owner, repo, number)
	if err != nil {
		log.Printf("error getting pull request: %v\n", err)
	}
	return resp, err
}

// GetPullRequestFiles list files from a pull request.
func (w *Worker) GetPullRequestFiles(owner string, repo string, number int, opts *v41.ListOptions) ([]*v41.CommitFile, error) {
	resp, _, err := w.Client.PullRequests.ListFiles(context.TODO(), owner, repo, number, opts)
	if err != nil {
		log.Printf("error getting files from pull request: %v\n", err)
	}
	return resp, err
}

// MergePullRequest merge a pull request.
func (w *Worker) MergePullRequest(owner string, repo string, number int, commitMessage string, options v41.PullRequestOptions) (*v41.PullRequestMergeResult, error) {
	resp, _, err := w.Client.PullRequests.Merge(context.TODO(), owner, repo, number, commitMessage, &options)
	if err != nil {
		log.Printf("error merging pull request: %v\n", err)
	}
	return resp, err
}

// PullRequestCreateReview create a pull request review.
func (w *Worker) PullRequestCreateReview(owner string, repo string, number int, review v41.PullRequestReviewRequest) (*v41.PullRequestReview, error) {
	resp, _, err := w.Client.PullRequests.CreateReview(context.TODO(), owner, repo, number, &review)
	if err != nil {
		log.Printf("error creating review: %v\n", err)
	}
	return resp, err
}

// IssueCreateComment create an issue or pull request comment.
func (w *Worker) IssueCreateComment(owner string, repo string, number int, comment v41.IssueComment) (*v41.IssueComment, error) {
	resp, _, err := w.Client.Issues.CreateComment(context.TODO(), owner, repo, number, &comment)
	if err != nil {
		log.Printf("error creating issue comment: %v\n", err)
	}
	return resp, err
}

// IssueUpdateComment update an issue or pull request comment.
func (w *Worker) IssueUpdateComment(owner string, repo string, commentID int, comment v41.IssueComment) (*v41.IssueComment, error) {
	resp, _, err := w.Client.Issues.EditComment(context.TODO(), owner, repo, int64(commentID), &comment)
	if err != nil {
		log.Printf("error updating issue comment: %v\n", err)
	}
	return resp, err
}

// IssueListComments list the comments of an issue or pull request.
func (w *Worker) IssueListComments(owner string, repo string, number int, opts *v41.IssueListCommentsOptions) ([]*v41.IssueComment, error) {
	resp, _, err := w.Client.Issues.ListComments(context.TODO(), owner, repo, number, opts)
	if err != nil {
		log.Printf("error getting issue list of comments: %v\n", err)
	}
	return resp, err
}

// GetCheckRun get a check run.
func (w *Worker) GetCheckRun(owner string, repo string, checkrunID int64) (*v41.CheckRun, error) {
	resp, _, err := w.Client.Checks.GetCheckRun(context.TODO(), owner, repo, checkrunID)
	if err != nil {
		log.Printf("error creating checkrun: %v\n", err)
	}
	return resp, err
}

// CreateCheckRun create a check run.
func (w *Worker) CreateCheckRun(owner string, repo string, checkrun v41.CreateCheckRunOptions) (*v41.CheckRun, error) {
	resp, _, err := w.Client.Checks.CreateCheckRun(context.TODO(), owner, repo, checkrun)
	if err != nil {
		log.Printf("error creating checkrun: %v\n", err)
	}
	return resp, err
}

// UpdateCheckRun update a check run.
func (w *Worker) UpdateCheckRun(owner string, repo string, checkrunID int64, checkRun v41.UpdateCheckRunOptions) (*v41.CheckRun, error) {
	resp, _, err := w.Client.Checks.UpdateCheckRun(context.TODO(), owner, repo, checkrunID, checkRun)
	if err != nil {
		log.Printf("error updating checkrun: %v\n", err)
	}
	return resp, err
}

// GetRef get a reference.
func (w *Worker) GetRef(owner string, repo string, ref string) (*v41.Reference, error) {
	resp, _, err := w.Client.Git.GetRef(context.TODO(), owner, repo, ref)
	if err != nil {
		log.Printf("error getting ref: %v\n", err)
	}
	return resp, err
}

// DeleteRef delete a reference.
func (w *Worker) DeleteRef(owner string, repo string, ref string) (*v41.Response, error) {
	resp, err := w.Client.Git.DeleteRef(context.TODO(), owner, repo, ref)
	if err != nil {
		log.Printf("error deleting ref: %v\n", err)
	}
	return resp, err
}

// CreatePulllRequestOverviewComment create a pull request comment with the report (overview diff).
func (w *Worker) CreatePulllRequestOverviewComment(owner *string, repo string, pullRequestNumber int) {
	if pullRequest, err := w.GetPullRequest(*owner, repo, pullRequestNumber); err == nil {
		currentIssueComment := w.GetPulllRequestOverviewComment(owner, repo, pullRequestNumber)
		extensions := make(map[string]int)
		if files, err := w.GetPullRequestFiles(*owner, repo, pullRequestNumber, nil); err == nil {
			for _, file := range files {
				ext := filepath.Ext(*file.Filename)
				value := extensions[ext]
				if value > 0 {
					extensions[ext] = value + 1
				} else {
					extensions[ext] = 1
				}
			}
		}
		var body bytes.Buffer
		body.WriteString(`<h3 align="center">Pull request Overview :checkered_flag:</h3>

---
<table align="center" border="10">
 <tr>
    <td>
`)
		body.WriteString(fmt.Sprintf("Commits &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &#8594; %d <br/>", *pullRequest.Commits))
		body.WriteString(fmt.Sprintf("Diff &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;&#8594; +%d  -%d <br/>", *pullRequest.Additions, *pullRequest.Deletions))
		body.WriteString(fmt.Sprintf("Changed Files &nbsp;&#8594; %d </td>", *pullRequest.ChangedFiles))
		body.WriteString(`<td>

| Ext | Amount |
`)
		body.WriteString(`|:-----:|--------|
`)
		for k, v := range extensions {
			body.WriteString(fmt.Sprintf(`| %s | %d |
`, k, v))
		}
		body.WriteString(`</td>
</tr>
</table>`)
		if currentIssueComment == nil {
			w.IssueCreateComment(*owner, repo, pullRequestNumber, v41.IssueComment{
				Body: v41.String(body.String()),
			})
		} else {
			w.IssueUpdateComment(*owner, repo, int(*currentIssueComment.ID), v41.IssueComment{
				Body: v41.String(body.String()),
			})
		}
	}
}

// GetPulllRequestOverviewComment get a pull request comment with the report (overview diff).
func (w *Worker) GetPulllRequestOverviewComment(owner *string, repo string, pullrequestID int) *v41.IssueComment {
	if comments, err := w.IssueListComments(*owner, repo, pullrequestID, nil); err == nil {
		for _, comment := range comments {
			if strings.HasPrefix(*comment.Body, "\u003ch3 align=\"center\"\u003ePull request Overview :checkered_flag:") {
				if strings.EqualFold("Bot", *comment.User.Type) && strings.Contains(*comment.User.AvatarURL, strconv.Itoa(w.Config.Github.AppID)) {
					return comment
				}
			}
		}
	}
	return nil
}
