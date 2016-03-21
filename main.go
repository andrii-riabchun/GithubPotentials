package main 

import (
    "github.com/go-martini/martini"
    "github.com/octokit/go-octokit/octokit"
    
    "github.com/andrewRyabchun/GithubTrendingPerspective/controllers"
)
var (
    token = octokit.TokenAuth{AccessToken: "kekek"}
)

func main() {
    
    
    m := martini.Classic()    
    
    client := octokit.NewClient(token)
    m.Map(&client)
    
    m.Get("/repoGraph/:criteria/:repopath/:timespan", controllers.GetGraph)
            
    m.Run()
}