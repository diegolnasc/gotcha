// Copyright 2021 Diego Lima. All rights reserved.

// Use of this source code is governed by a Apache license.
// license that can be found in the LICENSE file.
package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

// Settings represents the configuration aspects and authorization.
type Settings struct {
	Layout Layout `yaml:"layout"`
	Github Github `yaml:"github"`
}

// Layout represents permissions level and pull request funcionalities.
type Layout struct {
	Administration Administration `yaml:"administration"`
	PullRequest    PullRequest    `yaml:"pullRequest"`
}

// Administration represents general permissions.
type Administration struct {
	Permission Permission `yaml:"permission"`
}

// Permission represents high level user's permissions and repository permission.
type Permission struct {
	Users        []string       `yaml:"users"`
	Repositories []Repositories `yaml:"repositories"`
}

// Repositories represents low level permissions in repositories.
type Repositories struct {
	Repository Repository `yaml:"repository"`
}

// Repository represents user permission in a single repository.
type Repository struct {
	Name  string   `yaml:"name"`
	Users []string `yaml:"users"`
}

// PullRequest represents commands and funcionalities.
type PullRequest struct {
	EnableOverview        bool      `yaml:"enableOverview"`
	OverViewCommand       string    `yaml:"overViewCommand"`
	ApproveCommand        string    `yaml:"approveCommand"`
	RunTestSuiteCommand   string    `yaml:"runTestSuiteCommand"`
	MergeCommand          string    `yaml:"mergeCommand"`
	MergeAndDeleteCommand string    `yaml:"mergeAndDeleteCommand"`
	TestSuite             TestSuite `yaml:"testSuite"`
}

// TestSuite represents configuration for test cases.
type TestSuite struct {
	NamePattern string `yaml:"namePattern"`
	Reviewers   bool   `yaml:"reviewers"`
	Assignees   bool   `yaml:"assignees"`
	Labels      bool   `yaml:"labels"`
}

// Github represents github app owner configuration.
type Github struct {
	AppID              int      `yaml:"appId"`
	Organization       string   `yaml:"organization"`
	User               string   `yaml:"user"`
	WebhookSecret      string   `yaml:"webhookSecret"`
	PrivateKeyLocation string   `yaml:"privateKeyLocation"`
	PrivateKey         string   `yaml:"privateKey"`
	Events             []string `yaml:"events"`
}

// ReadConf initialize the yaml configuration
func (c *Settings) ReadConf() {
	yamlFile, err := ioutil.ReadFile("build/config-local.yaml")
	if err != nil {
		yamlFile, err = ioutil.ReadFile("build/config.yaml")
		if err != nil {
			log.Panic("yamlFile not found: ", err)
		}
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("error reading the yamlFile %v", err)
	}
}
