// Copyright 2021 Diego Lima. All rights reserved.

// Use of this source code is governed by a Apache license.
// license that can be found in the LICENSE file.
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/diegolnasc/gotcha/pkg/config"
	github "github.com/diegolnasc/gotcha/pkg/github"
)

type provider string

const (
	gitHub provider = "github"
)

func main() {
	p := flag.String("provider", "github", "Provider to run")
	flag.Parse()
	startProvider(provider(*p))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Panic(err)
	}
}

func startProvider(p provider) {
	switch p {
	case gitHub:
		config := &config.Settings{}
		config.ReadConf()
		worker := github.New(config)
		http.HandleFunc("/", worker.Handler)
	default:
		log.Panic("Provider not found")
	}
}
