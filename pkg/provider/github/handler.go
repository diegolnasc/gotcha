// Copyright 2021 Diego Lima. All rights reserved.

// Use of this source code is governed by a Apache license.
// license that can be found in the LICENSE file.

package github

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	ghwebhooks "github.com/go-playground/webhooks/v6/github"
)

// Events github events handler.
func (w *Worker) Events(c *gin.Context) {
	hook, err := ghwebhooks.New(ghwebhooks.Options.Secret(w.Config.Github.WebhookSecret))
	if err != nil {
		log.Panic(err)
	}
	var events []ghwebhooks.Event
	for _, event := range w.Config.Github.Events {
		events = append(events, ghwebhooks.Event(event))
	}
	payload, err := hook.Parse(c.Request, events...)
	if err != nil {
		if err == ghwebhooks.ErrEventNotFound {
			c.Writer.WriteHeader(http.StatusOK)
		} else {
			c.Writer.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	switch payload := payload.(type) {
	case ghwebhooks.PullRequestPayload:
		service := &PullRequestService{w: w}
		go service.processPullRequestEvent(&payload)
	case ghwebhooks.IssueCommentPayload:
		service := &IssueService{w: w}
		go service.processIssueCommentEvent(&payload)
	case ghwebhooks.CheckRunPayload:
		service := &CheckService{w: w}
		go service.processCheckRunEvent(&payload)
	default:
		log.Println("missing handler")
	}
	c.Writer.WriteHeader(http.StatusOK)
}
