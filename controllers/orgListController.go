package controllers

import (
    "net/http"
    "time"
    "sort"
    "strconv"
    "encoding/json"
    "github.com/go-martini/martini"
    "github.com/google/go-github/github" 
    "github.com/andrewRyabchun/GithubPotentials/models"
    "github.com/andrewRyabchun/GithubPotentials/helpers"
)

// GetOrgList GET /orgs/:criteria/:timespan
func GetOrgList(client *github.Client, params martini.Params) (int, []byte){
    days, err := strconv.Atoi(params["timespan"])
    if err != nil || days<1 {
        return http.StatusBadRequest, nil
    }
    date := time.Now().AddDate(0,0,-days)
       
    query:="repos:>1 type:org created:>"+helpers.FormatDate(date)
 
    result := models.OrgList{}
    switch params["criteria"] {
    case "total":
        orgs, err := listOrgByTotalCommits(client,query,date,days)
        if err != nil {
            println(err.Error())
            return http.StatusInternalServerError, nil
        }
        result.Criteria="total"
        result.Items=orgs
    case "avg":
        orgs, err := listOrgByAvgCommits(client,query,date,days)
        if err != nil {
            return http.StatusInternalServerError, nil
        }
        result.Criteria="avg"
        result.Items=orgs
    default:
        return http.StatusBadRequest, nil
    }
    resp, err := json.Marshal(result)
    if err != nil {
        return http.StatusInternalServerError, nil
    }
    return http.StatusOK, resp
}



func listOrgByTotalCommits(client *github.Client, query string, date time.Time, days int) ([]models.OrgEntry, error){
    var orgResult models.OrgSort
    i:=0
    opt:=&github.SearchOptions{
        Sort:   "followers",
        Order:  "asc",
        ListOptions: github.ListOptions{
            PerPage:10,
        },
    }  
    for {           
        i++
        orgs, resp, err := client.Search.Users(query, opt)
        if err!=nil{
            return nil, err
        }
       
        for _,org:= range orgs.Users{
            println(org.String())            
            repoCommits, err:=getTotalCommitsByOrg(client,*org.Login,date,days)
            if err != nil {
                return nil,err
            }
            
            repoCommitsFl := float32(repoCommits)
            result:=models.OrgEntry{
                Name: *org.Login,
                SortingCriteria: repoCommitsFl,
            }
            orgResult = append(orgResult, result)
            
        }
       
        if resp.NextPage == 0 || i==pagesCount {
            break
        }
        opt.Page = resp.NextPage
    }
    sort.Sort(orgResult)
    return orgResult, nil
}

func listOrgByAvgCommits(client *github.Client, query string, date time.Time, days int) ([]models.OrgEntry, error){
    var orgResult models.OrgSort
    i:=0
    opt:=&github.SearchOptions{
        Sort:   "followers",
        Order:  "asc",
        ListOptions: github.ListOptions{
            PerPage:10,
        },
    }
    for {           
        i++
        orgs, resp, err := client.Search.Users(query, opt)
        if err!=nil{
            return nil, err
        }
        
        for _,org:= range orgs.Users{
            
            repoCommits, err:=getTotalCommitsByOrg(client,*org.Login,date,days)
            
            avgCommits := float32(repoCommits)/float32(*org.Collaborators)
            if err != nil {
                return nil,err
            }

            result:=models.OrgEntry{
                Name: *org.Name,
                SortingCriteria: avgCommits,
            }
            orgResult = append(orgResult, result)
        }
       
        if resp.NextPage == 0 || i==pagesCount {
            break
        }
        opt.Page = resp.NextPage
    }
    sort.Sort(orgResult)
    return orgResult, nil
}

func reposByOrg(client *github.Client, org string) ([]github.Repository, error){
    opt:=&github.RepositoryListByOrgOptions{
        ListOptions: github.ListOptions{
            PerPage:100,
        }, 
    }
    var repos []github.Repository
    
    for {
        result, resp, err := client.Repositories.ListByOrg(org, opt)
        if err!=nil{
            return nil, err
        }
        repos = append(repos, result...)

       
        if resp.NextPage == 0 {
            break
        }
        opt.Page = resp.NextPage
    }
    return repos,nil
}

func getTotalCommitsByOrg(client *github.Client, org string, date time.Time, days int) (int, error){    
    repos,err:=reposByOrg(client,org)
    if err != nil {
        return -1, err
    }
    
    totalCommits := 0
    for _,repo:=range repos{       
        count,_,err := commits(client, *repo.Owner.Login, *repo.Name, date, days)
        if err != nil {
            return -1, err
        }
        totalCommits+=count
    }
    return totalCommits, nil
    
}