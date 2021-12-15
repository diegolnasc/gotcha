package github

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/diegolnasc/gotcha/pkg/config"
	v41 "github.com/google/go-github/v41/github"
)

type Worker struct {
	Config config.Settings
	Client v41.Client
}

func getOwner(auth *config.Settings) (*string, error) {
	if len(auth.Github.Organization) > 0 {
		return &auth.Github.Organization, nil
	} else if (len(auth.Github.User)) > 0 {
		return &auth.Github.User, nil
	}
	return nil, errors.New("owner not configured")
}

func New(auth *config.Settings) *Worker {
	var appTransport *ghinstallation.AppsTransport
	var installationTransport *ghinstallation.Transport
	var githubInstallation *v41.Installation
	var err error
	if len(auth.Github.PrivateKey) > 0 {
		appTransport, err = ghinstallation.NewAppsTransport(http.DefaultTransport, int64(auth.Github.AppID), []byte(auth.Github.PrivateKey))
	} else {
		appTransport, err = ghinstallation.NewAppsTransportKeyFromFile(http.DefaultTransport, int64(auth.Github.AppID), auth.Github.PrivateKeyLocation)
	}
	if err != nil {
		log.Panic("error initializing github authentication: ", err)
	}
	if len(auth.Github.Organization) > 0 {
		githubInstallation, _, err = v41.NewClient(&http.Client{Transport: appTransport}).Apps.FindOrganizationInstallation(context.TODO(), auth.Github.Organization)
	} else {
		githubInstallation, _, err = v41.NewClient(&http.Client{Transport: appTransport}).Apps.FindUserInstallation(context.TODO(), auth.Github.User)
	}
	if err != nil {
		log.Panic("error initializing github installation: ", err)
	}
	installationID := githubInstallation.GetID()
	installationTransport = ghinstallation.NewFromAppsTransport(appTransport, installationID)

	return &Worker{
		Config: *auth,
		Client: *v41.NewClient(&http.Client{Transport: installationTransport}),
	}
}
