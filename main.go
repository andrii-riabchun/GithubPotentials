package main

import (
    "fmt"
    "strconv"
    "log"
    "golang.org/x/oauth2"
	"github.com/go-martini/martini"
	"github.com/google/go-github/github"   
    "github.com/andrewRyabchun/GithubPotentials/helpers"
	"github.com/andrewRyabchun/GithubPotentials/controllers"
)

var app *martini.ClassicMartini
var conf *helpers.Configuration
const configFile = "config.json" 

func init(){
    app = martini.Classic()
    
    config,err :=helpers.LoadConfigFromFile(configFile)
    if err!=nil{
        fmt.Println("config error: ",err)
    }
    conf=config
    
    //map gh client to all controllers (DI)
    client := createAPIClient(config.GithubPersonalToken)
    
    app.Map(client)
    //routing
    app.Get("/repos/:criteria/:weeks", controllers.GetRepoList)
    app.Get("/orgs/:criteria/:weeks", controllers.GetOrgList)
    app.Get("/:owner/:repo/:days", controllers.GetRepoInfo)
    
    app.Use(func(c martini.Context, log *log.Logger, client *github.Client){
        c.Next()
        rate, _, err := client.RateLimits()
        if err != nil {
            return
        }
        
        log.Printf("Core rate: %d. \tReset: %s",rate.Core.Remaining, rate.Core.Reset)
        log.Printf("Search rate: %d. \tReset: %s",rate.Search.Remaining, rate.Search.Reset)
})
    
    
}

func main() {
	app.RunOnAddr(":"+strconv.Itoa(conf.Port))
}

//token - personal API token
func createAPIClient(token string) *github.Client{
  ts := oauth2.StaticTokenSource(
    &oauth2.Token{AccessToken: token},
  )
  tc := oauth2.NewClient(oauth2.NoContext, ts)
  return github.NewClient(tc)
}
