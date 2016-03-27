package controllers

import (
    "net/http"
    "time"
    "strconv"
    "encoding/json"
    "sort"
    "github.com/go-martini/martini"
    "github.com/google/go-github/github" 
    "github.com/andrewRyabchun/GithubPotentials/models"
    "github.com/andrewRyabchun/GithubPotentials/helpers"
)

// GetRepoList GET /repos/:criteria/:weeks
func GetRepoList(params martini.Params, req *http.Request) (int, []byte){
    weeks, err := strconv.Atoi(params["weeks"])
    if err != nil || weeks<1 || weeks>52 {
        return http.StatusBadRequest, nil
    }
    date := time.Now().AddDate(0,0,-weeks*7)
    days:=weeks*7
       
    query:="stars:>10 size:>1 pushed:>"+helpers.FormatDate(date)
 
    result := models.RepoList{}
    switch params["criteria"] {
    case "stars":
        repos, err := listByStars(query,date,days)
        if err != nil {
            return http.StatusInternalServerError, nil
        }
        result.Criteria="stars"
        result.Items=repos
    case "commits":
        repos, err := listByCommits(query,date,days)
        if err != nil {
            return http.StatusInternalServerError, nil
        }
        result.Criteria="commits"
        result.Items=repos
    case "contribs":
        repos, err := listByContribs(query,date,days)
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
    cache.Add(req.RequestURI, resp, expiration())
    return http.StatusOK, resp
}

func listByStars(query string, date time.Time, days int) ([]models.RepoEntry, error){
    var repos models.RepoSort
    result, err := searchRepos(query)
    if err != nil {
        return nil,err
    }
    for _,v:= range result{
        count,_,err := stars(*v.Owner.Login, *v.Name, date, days, true)
        if err != nil {
            return nil, err
        }
        if count==0{
            continue
        }
        repo := models.RepoEntry{
            FullName: *v.FullName,
            SortingCriteria: count,
        }
        repos = append(repos, repo)
    }

    sort.Sort(repos)
    return repos, nil
}

func listByCommits(query string, date time.Time, days int) ([]models.RepoEntry, error){
    var repos models.RepoSort
    
    result, err := searchRepos(query)
    if err != nil {
        return nil,err
    }
    for _,v:= range result{
        count,_,err := commits(*v.Owner.Login, *v.Name, date, days)
        if err != nil {
            return nil, err
        }
        if count==0{
            continue
        }
        repo := models.RepoEntry{
            FullName: *v.FullName,
            SortingCriteria: count,
        }
        repos = append(repos, repo)
    }
    sort.Sort(repos)
    return repos, nil
}

func listByContribs(query string, date time.Time, days int) ([]models.RepoEntry, error){
    var repos models.RepoSort
    result, err := searchRepos(query)
    if err != nil {
        return nil,err
    }
    for _,v:= range result{
        count,_,err := contributors(*v.Owner.Login, *v.Name, date, days)
        if err != nil {
            return nil, err
        }
        if count==0{
            continue   
        }
        repo := models.RepoEntry{
            FullName: *v.FullName,
            SortingCriteria: count,
        }
        repos = append(repos, repo)
    }
       
    sort.Sort(repos)
    return repos, nil
}

func searchRepos(query string) ([]github.Repository,error){
    if cached,is:=cache.Get("reposearch:"+query);is{
        return cached.([]github.Repository), nil
    }
    i:=0
    opt:=&github.SearchOptions{
        Sort:   "stars",
        Order:  "asc",
        ListOptions: github.ListOptions{
            PerPage:100,
        },
    }
    
    var repos []github.Repository
    
    for {           
        i++
        result, resp, err := client.Search.Repositories(query, opt)
        if err!=nil{
            return nil, err
        }
        repos = append(repos,result.Repositories...)
       
        if resp.NextPage == 0 || i==repoPagesCount {
            break
        }
        opt.Page = resp.NextPage
    }
    cache.Add("reposearch:"+query, repos, expiration())
    return repos, nil
}