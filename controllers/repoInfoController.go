package controllers

import (
    "net/http"
    "time"
    "strconv"
    "encoding/json"
    "github.com/go-martini/martini"
    "github.com/google/go-github/github" 
    "github.com/andrewRyabchun/GithubTrendingPerspective/models"
    "github.com/andrewRyabchun/GithubTrendingPerspective/helpers"
)
var gh *github.Client 
   
//GetStarsInfo - GET /repoinfo/stars/:owner/:repo/:timespan
func GetStarsInfo(client *github.Client, params martini.Params) (int, []byte) {
    owner:=params["owner"]
    repo :=params["repo"]   
    days, err := strconv.Atoi(params["timespan"])
    if err != nil {
        //cant parse days count
        return http.StatusBadRequest, nil
    }
    date := time.Now().AddDate(0,0,-days)
      
    opt := &github.ListOptions{PerPage:100}
    
    daysStarsDict := make(map[int]int,days)
    totalStars := 0
    for {
        //no timestamps =(, going to contribute go-github
        stargazers, resp, err := helpers.ListStargazers(client,owner,repo,opt)
        if err!=nil{
            return http.StatusInternalServerError,nil
        }
        
        for _, sg := range stargazers{
            if sg.StarredAt.Time.After(date){                 
                dayBeforeNow := helpers.DaysSinceNow(sg.StarredAt.Time)
                daysStarsDict[dayBeforeNow]++     
                totalStars++                
            }
        }
        
        if resp.NextPage == 0 {
            break
        }
	    opt.Page = resp.NextPage
    }
    resArr := make([]int, days)
    
    for i:=0;i<days;i++{
        resArr[days-1-i]=daysStarsDict[i]
    }
    
    stars := &models.RepoStars{
        Owner: owner,
        Repo: repo,
        TotalStars: totalStars,
        DayCount: days,
        DayStars: resArr,
    }
    
    resp, err := json.Marshal(stars)
    if err != nil {
        return http.StatusInternalServerError, nil
    }
    return http.StatusOK, resp
}
//GetCommitsInfo GET /repoinfo/commits/:owner/:repo/:timespan
func GetCommitsInfo(client *github.Client, params martini.Params) (int, []byte) {
    owner:=params["owner"]
    repo :=params["repo"]   
    days, err := strconv.Atoi(params["timespan"])
    if err != nil {
        return http.StatusBadRequest, nil
    }
    date := time.Now().AddDate(0,0,-days)
      
    opt := &github.CommitsListOptions{Since:date}
    
    daysCommitsDict := make(map[int]int,days)
    totalCommits := 0
    for {
        commits, resp, err := client.Repositories.ListCommits(owner, repo, opt)
        if err!=nil{
            return http.StatusInternalServerError,nil
        }
        
        for _, commit := range commits{              
                dayBeforeNow := helpers.DaysSinceNow(*commit.Commit.Committer.Date)
                daysCommitsDict[dayBeforeNow]++     
                totalCommits++                
        }
        
        if resp.NextPage == 0 {
            break
        }
	    opt.Page = resp.NextPage
    }
    resArr := make([]int, days)
    
    for i:=0;i<days;i++{
        resArr[days-1-i]=daysCommitsDict[i]
    }
    
    stars := &models.RepoCommits{
        Owner: owner,
        Repo: repo,
        TotalCommits: totalCommits,
        DayCount: days,
        DayCommits: resArr,
    }
    
    resp, err := json.Marshal(stars)
    if err != nil {
        return http.StatusInternalServerError, nil
    }
    return http.StatusOK, resp
}
//GetContributorsInfo GET /repoinfo/contributors/:owner/:repo/:timespan
func GetContributorsInfo(client *github.Client, params martini.Params) (int, []byte) {
    owner:=params["owner"]
    repo :=params["repo"]   
    days, err := strconv.Atoi(params["timespan"])
    if err != nil {
        return http.StatusBadRequest, nil
    }
    date := time.Now().AddDate(0,0,-days)
      
    opt := &github.CommitsListOptions{Since:date}
        
    daysContribsDict := make(map[int][]int,days) //[day: [ids..]...]
    uniqueContribs := make(map[int]interface{}, 31)
    for {
        commits, resp, err := client.Repositories.ListCommits(owner, repo, opt)
        if err!=nil{
            return http.StatusInternalServerError,nil
        }
        for _, commit := range commits{                        
            dayBeforeNow := helpers.DaysSinceNow(*commit.Commit.Committer.Date)
            
            if commit.Author!=nil{
                daysContribsDict[dayBeforeNow]=append(daysContribsDict[dayBeforeNow],*commit.Author.ID)
                uniqueContribs[*commit.Author.ID]=nil
            }else if commit.Committer !=nil{
                daysContribsDict[dayBeforeNow]=append(daysContribsDict[dayBeforeNow],*commit.Committer.ID)
                uniqueContribs[*commit.Committer.ID]=nil //add to unique collection                         
            }else{
                continue
            }
                
        }
        
        if resp.NextPage == 0 {
            break
        }
	    opt.Page = resp.NextPage
    }
    resArr := make([]int, days)
    
    for i:=0;i<days;i++{
        daysContribsDict[i] = helpers.RemoveDuplicates(daysContribsDict[i])
        resArr[days-1-i]=len(daysContribsDict[i])
    }
    
    stars := &models.RepoContributors{
        Owner: owner,
        Repo: repo,      
        DayCount: days,
        DayContribs: resArr,
        UniqueContribs: len(uniqueContribs),
    }
    
    resp, err := json.Marshal(stars)
    if err != nil {
        return http.StatusInternalServerError, nil
    }
    return http.StatusOK, resp
}