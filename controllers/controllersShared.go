package controllers

import (
    "time"
    "github.com/google/go-github/github"
    "github.com/facebookgo/inmem"   
)

const (
    repoPagesCount = 3
    orgPagesCount = 2
)

var client *github.Client
var cache inmem.Cache
const expirationMinutes = 60

// Init services
func Init(cacheArg *inmem.Cache, clientArg *github.Client){
    cache=*cacheArg
    client=clientArg
}

func expiration() time.Time{
    return time.Now().Add(time.Minute*expirationMinutes)
}

func reposByUser(user string) ([]github.Repository, error){
    if cached,is:=cache.Get("repos:"+user);is{
        return cached.([]github.Repository), nil
    }
    
    opt:=&github.RepositoryListOptions{
        ListOptions: github.ListOptions{
            PerPage:100,
        }, 
    }
    var repos []github.Repository
    
    for {         
        result, resp, err := client.Repositories.List(user,opt)
        if err!=nil{
            return nil, err
        }
        repos = append(repos, result...)
     
        if resp.NextPage == 0 {
            break
        }
        opt.Page = resp.NextPage
    }
    cache.Add("repos:"+user, repos, expiration())
    return repos,nil
}