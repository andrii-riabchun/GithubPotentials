package controllers

import (
    "net/http"
    "time"
    //"sync"
    "strconv"
    "encoding/json"
    "sort"
    "github.com/go-martini/martini"
    "github.com/google/go-github/github" 
    "github.com/andrewRyabchun/GithubTrendingPerspective/models"
    "github.com/andrewRyabchun/GithubTrendingPerspective/helpers"
)
const pagesCount = 1
// GetRepoList GET /repos/:criteria/:timespan
func GetRepoList(client *github.Client, params martini.Params) (int, []byte){
    days, err := strconv.Atoi(params["timespan"])
    if err != nil || days<1 {
        return http.StatusBadRequest, nil
    }
    date := time.Now().AddDate(0,0,-days)
       
    query:="stars:>10 size:>1 pushed:>"+helpers.FormatDate(date)
 
    result := models.RepoList{}
    switch params["criteria"] {
    case "stars":
        repos, err := listByStars(client,query,date,days)
        if err != nil {
            return http.StatusInternalServerError, nil
        }
        result.Criteria="stars"
        result.Items=repos
    case "commits":
        repos, err := listByCommits(client,query,date,days)
        if err != nil {
            return http.StatusInternalServerError, nil
        }
        result.Criteria="commits"
        result.Items=repos
    case "contribs":
        repos, err := listByContribs(client,query,date,days)
        if err != nil {
            return http.StatusInternalServerError, nil
        }
        result.Criteria="contribs"
        result.Items=repos
    default:
        return http.StatusBadRequest, nil
    }
    resp, err := json.Marshal(result)
    if err != nil {
        return http.StatusInternalServerError, nil
    }
    return http.StatusOK, resp
}

func listByStars(client *github.Client, query string, date time.Time, days int) ([]models.RepoEntry, error){
    i:=0
    opt:=&github.SearchOptions{
        Sort:   "stars",
        Order:  "asc",
        ListOptions: github.ListOptions{
            PerPage:100,
        },
    }
    var repos models.RepoSort
    
    for {           
        i++
        result, resp, err := client.Search.Repositories(query, opt)
        if err!=nil{
            return nil, err
        }
        
        for _,v:= range result.Repositories{
            count,_,err := stars(client, *v.Owner.Login, *v.Name, date, days)
            if err != nil {
                return nil, err
            }
            repo := models.RepoEntry{
                FullName: *v.FullName,
                SortingCriteria: count,
            }
            repos = append(repos, repo)
        }
       
        if resp.NextPage == 0 || i==pagesCount {
            break
        }
        opt.Page = resp.NextPage
    }
    sort.Sort(repos)
    return repos, nil
}

func listByCommits(client *github.Client, query string, date time.Time, days int) ([]models.RepoEntry, error){
    i:=0
    opt:=&github.SearchOptions{
        Sort:   "stars",
        Order:  "asc",
        ListOptions: github.ListOptions{
            PerPage:100,
        },
    }
    var repos models.RepoSort
    
    for {           
        i++
        result, resp, err := client.Search.Repositories(query, opt)
        if err!=nil{
            return nil, err
        }
        println(i)
        for _,v:= range result.Repositories{
            count,_,err := commits(client, *v.Owner.Login, *v.Name, date, days)
            if err != nil {
                return nil, err
            }
            repo := models.RepoEntry{
                FullName: *v.FullName,
                SortingCriteria: count,
            }
            repos = append(repos, repo)
        }
       
        if resp.NextPage == 0 || i==pagesCount {
            break
        }
        opt.Page = resp.NextPage
    }
    sort.Sort(repos)
    return repos, nil
}

func listByContribs(client *github.Client, query string, date time.Time, days int) ([]models.RepoEntry, error){
    i:=0
    opt:=&github.SearchOptions{
        Sort:   "stars",
        Order:  "asc",
        ListOptions: github.ListOptions{
            PerPage:100,
        },
    }
    var repos models.RepoSort
    
    for {           
        i++
        result, resp, err := client.Search.Repositories(query, opt)
        if err!=nil{
            return nil, err
        }
        
        for _,v:= range result.Repositories{
            count,_,err := contributors(client, *v.Owner.Login, *v.Name, date, days)
            if err != nil {
                return nil, err
            }
            repo := models.RepoEntry{
                FullName: *v.FullName,
                SortingCriteria: count,
            }
            repos = append(repos, repo)
        }
       
        if resp.NextPage == 0 || i==pagesCount {
            break
        }
        opt.Page = resp.NextPage
    }
    sort.Sort(repos)
    return repos, nil
}