package main

import (
	"encoding/json"
	"io"

	"github.com/artisresistance/githubpotentials/github"
	"time"
)

type config struct {
	Token           string
	OutputPath      string
	OutCount        int
	FetchPagesCount int
}

func loadConfig(r io.ReadCloser) (config, error) {
	defer r.Close()
	conf := config{}

	err := json.NewDecoder(r).Decode(&conf)

	return conf, err
}

type result struct {
	Metadata       meta
	ByCommits      []github.Repository
	ByStars        []github.Repository
	ByContributors []github.Repository
}

type meta struct {
	Updated          time.Time
	APICallsRemained int
	Errors           int
	DurationSec      int
	Reset            time.Time
}

func (r result) Write(wc io.WriteCloser) error {
	err := json.NewEncoder(wc).Encode(r)
	return err
}
