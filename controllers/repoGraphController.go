package controllers

import (
    "github.com/go-martini/martini"
    "github.com/google/go-github/github" 
    "net/http"
)
    
//GetGraph - GET /repoGraph/:criteria/:repopath/:timespan
func GetGraph(client *github.Client, p *martini.Params) (int, string) {
    return http.StatusOK, "kek"
}