package main

import (
	"log"
	"net/http"

	config "github.com/diegolnasc/gotcha/pkg/config"
	github "github.com/diegolnasc/gotcha/pkg/github"
)

type Provider string

const (
	GitHub Provider = "github"
)

func main() {
	startProvider(Provider("github"))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Panic(err)
	}
}

func startProvider(provider Provider) {
	switch provider {
	case GitHub:
		config := &config.Settings{}
		config.ReadConf()
		worker := github.New(config)
		http.HandleFunc("/", worker.Handler)
	default:
		log.Panic("Provider not found")
	}
}
