package controllers

import (
    "github.com/octokit/go-octokit/octokit"
    "github.com/go-martini/martini"
    "net/http")
    
//GetGraph kek
func GetGraph(github *octokit.Client, p *martini.Params) (int, string) {
    return http.StatusOK, "kek"
}