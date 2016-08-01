package main

import (
    "io"
    "encoding/json"

    potentials "github.com/artisresistance/githubpotentials"
)

type Config struct{
    Token string
    OutCount int
}

func LoadConfig(r io.ReadCloser) (Config, error) {
    defer r.Close()
    conf := Config{}

    err := json.NewDecoder(r).Decode(&conf)

    return conf, err
}

type PotentialsResult struct{
    Updated int64
    Fetched int
    SortedBy string
    Errors int
    Items []potentials.Repository
}


func (r PotentialsResult) Write(wc io.WriteCloser) error {
    err := json.NewEncoder(wc).Encode(r)
    return err
}