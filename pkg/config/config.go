package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Settings struct {
	Layout Layout `yaml:"layout"`
	Github Github `yaml:"github"`
}

type Layout struct {
	Administration Administration `yaml:"administration"`
	PullRequest    PullRequest    `yaml:"pullRequest"`
}

type Administration struct {
	Permission Permission `yaml:"permission"`
}

type Permission struct {
	Users        []string       `yaml:"users"`
	Repositories []Repositories `yaml:"repositories"`
}

type Repositories struct {
	Repository Repository `yaml:"repository"`
}

type Repository struct {
	Name  string   `yaml:"name"`
	Users []string `yaml:"users"`
}

type PullRequest struct {
	ApproveCommand        string    `yaml:"approveCommand"`
	ReRunTestSuiteCommand string    `yaml:"reRunTestSuiteCommand"`
	TestSuite             TestSuite `yaml:"testSuite"`
}

type TestSuite struct {
	NamePattern string `yaml:"namePattern"`
	Reviewers   bool   `yaml:"reviewers"`
	Assignees   bool   `yaml:"assignees"`
	Labels      bool   `yaml:"labels"`
}

type Github struct {
	AppID              int      `yaml:"appId"`
	Organization       string   `yaml:"organization"`
	User               string   `yaml:"user"`
	WebhookSecret      string   `yaml:"webhookSecret"`
	PrivateKeyLocation string   `yaml:"privateKeyLocation"`
	PrivateKey         string   `yaml:"privateKey"`
	Events             []string `yaml:"events"`
}

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
