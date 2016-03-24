package main

import (
    "fmt"
    "strconv"
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
    
    config,err :=helpers.LoadConfigFromFile(configFile)
    if err!=nil{
        fmt.Println("config error: ",err)
    }
    conf=config
    
    //map gh client to all controllers (DI)
    client := createAPIClient(config.GithubPersonalToken)
    
    app.Map(client)
    //routing
    app.Get("/:owner/:repo/:timespan", controllers.GetRepoInfo)
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
