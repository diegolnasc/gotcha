// Copyright 2021 Diego Lima. All rights reserved.

// Use of this source code is governed by a Apache license.
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/diegolnasc/gotcha/pkg/model"
	v41 "github.com/google/go-github/v41/github"
)

type service struct {
	w *GitHubWorker
}

// Worker represents the general configuration to run the github provider.
type GitHubWorker struct {
	Config model.Settings
	Client v41.Client
}

// getOwner get an owner of a provider.
func getOwner(auth *model.Settings) (*string, error) {
	if len(auth.Github.Organization) > 0 {
		return &auth.Github.Organization, nil
	} else if (len(auth.Github.User)) > 0 {
		return &auth.Github.User, nil
	}
	return nil, errors.New("owner not configured")
}

// isUserAuthorized check if the user running a command is authorized.
func isUserAuthorized(auth *model.Settings, user *string, repo *string) bool {
	result := false
	if user != nil && repo != nil {
		for _, u := range auth.Layout.Administration.Permission.Users {
			if u == *user {
				result = true
				break
			}
		}
		if !result {
			for _, r := range auth.Layout.Administration.Permission.Repositories {
				if r.Repository.Name == *repo {
					for _, u := range r.Repository.Users {
						if u == *user {
							result = true
							break
						}
					}
				}
			}
		}
	}
	return result
}

// Initialize a GitHub handler Worker.
func New(cfg *model.Settings) *GitHubWorker {
	var appTransport *ghinstallation.AppsTransport
	var installationTransport *ghinstallation.Transport
	var githubInstallation *v41.Installation
	var err error
	if len(cfg.Github.PrivateKey) > 0 {
		appTransport, err = ghinstallation.NewAppsTransport(http.DefaultTransport, int64(cfg.Github.AppID), []byte(cfg.Github.PrivateKey))
	} else {
		appTransport, err = ghinstallation.NewAppsTransportKeyFromFile(http.DefaultTransport, int64(cfg.Github.AppID), cfg.Github.PrivateKeyLocation)
	}
	if err != nil {
		log.Panic("error initializing github authentication: ", err)
	}
	if len(cfg.Github.Organization) > 0 {
		githubInstallation, _, err = v41.NewClient(&http.Client{Transport: appTransport}).Apps.FindOrganizationInstallation(context.TODO(), cfg.Github.Organization)
	} else {
		githubInstallation, _, err = v41.NewClient(&http.Client{Transport: appTransport}).Apps.FindUserInstallation(context.TODO(), cfg.Github.User)
	}
	if err != nil {
		log.Panic("error initializing github installation: ", err)
	}
	installationID := githubInstallation.GetID()
	installationTransport = ghinstallation.NewFromAppsTransport(appTransport, installationID)

	return &GitHubWorker{
		Config: *cfg,
		Client: *v41.NewClient(&http.Client{Transport: installationTransport}),
	}
}
