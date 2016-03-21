package main

import (
    "fmt"
    "golang.org/x/oauth2"
	"github.com/go-martini/martini"
	"github.com/google/go-github/github"   
    "github.com/andrewRyabchun/GithubTrendingPerspective/helpers"
	"github.com/andrewRyabchun/GithubTrendingPerspective/controllers"
)

var app *martini.ClassicMartini
var conf *helpers.Configuration
const configFile = "config.json" 

func init(){
    app = martini.Classic()
    //routing
    app.Get("/repoGraph/:criteria/:repopath/:timespan", controllers.GetGraph)
    
   
    
    conf,err:=helpers.LoadConfigFromFile(configFile)
    if err!=nil{
        fmt.Println("config error: ",err)
    }
    
    //map gh client to all controllers (DI)
    client:=createAPIClient(conf.GithubPersonalToken)
    app.Map(&client)
        
}

func main() {
	app.RunOnAddr(":"+string(conf.Port))
}

func createAPIClient(token string) *github.Client{
  ts := oauth2.StaticTokenSource(
    &oauth2.Token{AccessToken: token},
  )
  tc := oauth2.NewClient(oauth2.NoContext, ts)
  return github.NewClient(tc)
}
