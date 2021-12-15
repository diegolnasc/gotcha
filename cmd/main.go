package main

import (
	"log"
	"net/http"

	config "github.com/diegolnasc/gotcha/pkg/config"
	github "github.com/diegolnasc/gotcha/pkg/github"
)

func main() {
	config := &config.Settings{}
	config.GetConf()
	worker := initGitub(config)
	http.HandleFunc("/", worker.Handler)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Panic(err)
	}
}

func initGitub(config *config.Settings) *github.Worker {
	return github.New(config)
}