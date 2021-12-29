package main

import (
	"flag"
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
	provider := flag.String("provider", "github", "Provider to run")
	flag.Parse()
	startProvider(Provider(*provider))
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
