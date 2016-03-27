package main

import (
    "fmt"
    "strconv"
    "log"
    "net/http"
    "golang.org/x/oauth2"
    "github.com/facebookgo/inmem"
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
    
    //map caching
    cache:=inmem.NewLocked(1000)
    app.Map(&cache)
    
    controllers.Init(&cache,client)

    //routing
    app.Get("/repos/:criteria/:weeks", controllers.GetRepoList)
    app.Get("/orgs/:criteria/:weeks", controllers.GetOrgList)
    app.Get("/:owner/:repo/:days", controllers.GetRepoInfo)
    
    //init middleware
    app.Use(cachingResponse)
    app.Use(rateLimitsMonitor)
      
    
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

// middlewares
func cachingResponse(c martini.Context, log *log.Logger, cache *inmem.Cache, req *http.Request, resp http.ResponseWriter){
    if bytes, cached := (*cache).Get(req.RequestURI);cached{
        resp.Write(bytes.([]byte))
    }else{
        c.Next()
    }
    log.Printf("Items in cache: %d",(*cache).Len())   
}

func rateLimitsMonitor(c martini.Context, log *log.Logger, client *github.Client, resp http.ResponseWriter){
    before,_,err:=client.RateLimits()
    if err != nil {
        log.Println(err.Error())
    }
    
    if (before.Core.Remaining==0 || before.Search.Remaining==0){
        var latest string
        if before.Core.Reset.Time.After(before.Search.Reset.Time){
            latest= before.Core.Reset.String()
        }else{
            latest= before.Search.Reset.String()
        }
        
        log.Println("API unavaiable. Come back "+latest)
        resp.WriteHeader(http.StatusTeapot)
    }else{
        c.Next()
        
        after, _, err := client.RateLimits()
        if err != nil {
            log.Println(err.Error())
        }
        
        log.Printf("Core rate remained: %d. \tUsed: %d. \tReset: %s",
                after.Core.Remaining,
                before.Core.Remaining-after.Core.Remaining, 
                after.Core.Reset)
        
        log.Printf("Search rate remained: %d. \tUsed: %d. \tReset: %s",
                after.Search.Remaining,
                before.Search.Remaining-after.Search.Remaining, 
                after.Search.Reset)
        }
}