// Copyright 2021 Diego Lima. All rights reserved.

// Use of this source code is governed by a Apache license.
// license that can be found in the LICENSE file.
package github

import (
	"log"
	"net/http"

	ghwebhooks "github.com/go-playground/webhooks/v6/github"
)

// Handler webhook event handler.
func (w *Worker) Handler(response http.ResponseWriter, request *http.Request) {
	log.Println(request.Header.Get("X-GitHub-Event"))
	hook, err := ghwebhooks.New(ghwebhooks.Options.Secret(w.Config.Github.WebhookSecret))
	if err != nil {
		log.Panic(err)
	}
	var events []ghwebhooks.Event
	for _, event := range w.Config.Github.Events {
		events = append(events, ghwebhooks.Event(event))
	}
	payload, err := hook.Parse(request, events...)
	if err != nil {
		if err == ghwebhooks.ErrEventNotFound {
			response.WriteHeader(http.StatusOK)
		} else {
			response.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	switch payload := payload.(type) {
	case ghwebhooks.PullRequestPayload:
		service := &PullRequestService{w: *w}
		go service.processPullRequestEvent(&payload)
	case ghwebhooks.IssueCommentPayload:
		service := &IssueService{w: *w}
		go service.processIssueCommentEvent(&payload)
	case ghwebhooks.CheckRunPayload:
		service := &CheckService{w: *w}
		go service.processCheckRunEvent(&payload)
	default:
		log.Println("missing handler")
	}

	response.WriteHeader(http.StatusOK)
}
