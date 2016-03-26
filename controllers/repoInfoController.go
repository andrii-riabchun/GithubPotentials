package controllers

import (
    "net/http"
    "time"
    "sync"
    "strconv"
    "encoding/json"
    "github.com/go-martini/martini"
    "github.com/google/go-github/github" 
    "github.com/andrewRyabchun/GithubPotentials/models"
    "github.com/andrewRyabchun/GithubPotentials/helpers"
)
   
// GetRepoInfo GET /:owner/:repo/:days
func GetRepoInfo(client *github.Client, params martini.Params) (int, []byte){
    owner:=params["owner"]
    repo :=params["repo"]   
    days, err := strconv.Atoi(params["days"])
    if err != nil || days<1 {
        return http.StatusBadRequest, nil
    }
    date := time.Now().AddDate(0,0,-days)
    
    joiner:=new(sync.WaitGroup)
    joiner.Add(3) 
    
    var isError bool
    var starsCount int
    var starsData []int    
    go func() {
        defer joiner.Done()
        var starsErr error
        starsCount,starsData,starsErr=stars(client, owner, repo, date, days, false)
        if starsErr!=nil{
            isError=true
        }
    }()
    
    var contribCount int
    var contribData []int
    go func(){
        defer joiner.Done()
        var contribErr error
        contribCount, contribData, contribErr = contributors(client, owner, repo, date, days)
        if contribErr!=nil{
            isError=true
        }
    }()
    
    var commitsCount int
    var commitsData []int
    go func(){
        defer joiner.Done()
        var commitErr error
        commitsCount, commitsData, commitErr = commits(client, owner, repo, date, days)
        if commitErr!=nil{
            isError=true
        }
    }()
    
    joiner.Wait()
    if isError {
        return http.StatusInternalServerError, nil
    }
      
    result:=&models.RepoInfo{
        Owner:          owner,
        Repo:           repo,
        DayCount:       days,
        Stars:          starsCount,
        StarsData:      starsData,
        Commits:        commitsCount,
        CommitsData:    commitsData,        
        Contribs:       contribCount,
        ContribsData:   contribData,
        
    }
    resp, err := json.Marshal(result)
    if err != nil {
        return http.StatusInternalServerError, nil
    }
    return http.StatusOK, resp
}

func contributors(client *github.Client,owner,repo string, date time.Time, days int) (int, []int, error){
    opt := &github.CommitsListOptions{Since:date}        
    daysContribsDict := make(map[int][]int,days) //[day: [ids..]...]
    uniqueContribs := make(map[int]interface{}, days)
    for {
        commits, resp, err := client.Repositories.ListCommits(owner, repo, opt)
        if err!=nil{
            return 0,nil,err
        }
        for _, commit := range commits{                        
            daySinceNow := helpers.DaysSinceNow(*commit.Commit.Committer.Date)
            
            if commit.Author!=nil{
                daysContribsDict[daySinceNow]=append(daysContribsDict[daySinceNow],*commit.Author.ID)
                uniqueContribs[*commit.Author.ID]=nil
            }else if commit.Committer !=nil{
                daysContribsDict[daySinceNow]=append(daysContribsDict[daySinceNow],*commit.Committer.ID)
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
    return len(uniqueContribs), resArr, nil
}

func commits(client *github.Client,owner,repo string, date time.Time, days int) (int, []int, error){
    opt := &github.CommitsListOptions{Since:date}
    
    daysCommitsDict := make(map[int]int,days)
    totalCommits := 0
    for {
        
        commits, resp, err := client.Repositories.ListCommits(owner, repo, opt)
        if (len(commits)) == 0{
            continue
        }
            
        if err!=nil{
            return 0,nil,err
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
    return totalCommits, resArr, nil
}

func stars(client *github.Client,owner,repo string, date time.Time, days int, omitPopular bool) (int, []int, error){
    opt := &github.ListOptions{PerPage:100}
    
    daysStarsDict := make(map[int]int,days)
    totalStars := 0
    for {
        //no timestamps =(, going to contribute go-github
        stargazers, resp, err := helpers.ListStargazers(client,owner,repo,opt)
        if err!=nil{
            return 0,nil,err
        }
        if omitPopular && len(stargazers)>100{
            continue
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
    return totalStars, resArr, nil
}