package main

import (
	"encoding/json"
	"io"

	potentials "github.com/artisresistance/githubpotentials"
)

type config struct {
	Token    string
	OutputPath string
	OutCount int
	FetchPagesCount int
}

func loadConfig(r io.ReadCloser) (config, error) {
	defer r.Close()
	conf := config{}

	err := json.NewDecoder(r).Decode(&conf)

	return conf, err
}

type potentialsResult struct {
	Updated  int64
	Fetched  int
	SortedBy string
	Errors   int
	Items    []potentials.Repository
}

func (r potentialsResult) Write(wc io.WriteCloser) error {
	err := json.NewEncoder(wc).Encode(r)
	return err
}
