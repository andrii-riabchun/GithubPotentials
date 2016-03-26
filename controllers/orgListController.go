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
const orgPagesCount = 3
// GetOrgList GET /orgs/:criteria/:weeks
func GetOrgList(client *github.Client, params martini.Params) (int, []byte){
    weeks, err := strconv.Atoi(params["weeks"])
    if err != nil || weeks<1 || weeks>52 {
        return http.StatusBadRequest, nil
    }
    
    date := time.Now().AddDate(0,0,-weeks*7)      
    query:="repos:>1 type:org created:>"+helpers.FormatDate(date)
 
    result := models.OrgList{}
    switch params["criteria"] {
    case "total":
        orgs, err := listOrgByTotalCommits(client,query,weeks)
        if err != nil {
            return http.StatusInternalServerError, nil
        }
        result.Criteria="total"
        result.Items=orgs
    case "avg":
        orgs, err := listOrgByAvgCommits(client,query,weeks)
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



func listOrgByTotalCommits(client *github.Client, query string, weeks int) ([]models.OrgEntry, error){
    var orgResult models.OrgSort
    i:=0
    opt:=&github.SearchOptions{
        Sort:   "followers",
        Order:  "asc",
        ListOptions: github.ListOptions{
            PerPage:100,
        },
    }  
    for {           
        i++
        orgs, resp, err := client.Search.Users(query, opt)
        if err!=nil{
            return nil, err
        }
       
        for _,org:= range orgs.Users{                   
            repoCommits, err:=getTotalCommitsByOrg(client,*org.Login,weeks)
            if err != nil {
                return nil,err
            }
            if repoCommits==0 {
                continue
            }           
            repoCommitsFl := float32(repoCommits)
            result:=models.OrgEntry{
                Name: *org.Login,
                SortingCriteria: repoCommitsFl,
            }
            orgResult = append(orgResult, result)
            
        }
       
        if resp.NextPage == 0 || i==orgPagesCount {
            break
        }
        opt.Page = resp.NextPage
    }
    sort.Sort(orgResult)
    return orgResult, nil
}

func listOrgByAvgCommits(client *github.Client, query string, weeks int) ([]models.OrgEntry, error){
    var orgResult models.OrgSort
    i:=0
    opt:=&github.SearchOptions{
        Sort:   "followers",
        Order:  "asc",
        ListOptions: github.ListOptions{
            PerPage:100,
        },
    }
    for {           
        i++
        orgs, resp, err := client.Search.Users(query, opt)
        if err!=nil{
            return nil, err
        }
        println(len(orgs.Users)) 
        for _,orgUser:= range orgs.Users{
                      
            repoCommits, err:=getTotalCommitsByOrg(client,*orgUser.Login,weeks)
            if repoCommits==0 {
                continue
            }    
            org, _, err := client.Organizations.Get(*orgUser.Login)
            if err != nil || org.Collaborators==nil || *org.Collaborators==0{
                continue
            } 
            members, err := orgMembersCount(client,*orgUser.Login)
            if err != nil {
                return nil, err
            }
            if members==0{
                continue
            }
            avgCommits := float32(repoCommits)/float32(members)
            if err != nil {
                return nil,err
            }
            
            result:=models.OrgEntry{
                Name: *orgUser.Login,
                SortingCriteria: avgCommits,
            }
            
            println(result.SortingCriteria)
            
            orgResult = append(orgResult, result)
        }
       
        if resp.NextPage == 0 || i==orgPagesCount {
            break
        }
        opt.Page = resp.NextPage
    }
    sort.Sort(orgResult)
    return orgResult, nil
}

func getTotalCommitsByOrg(client *github.Client, org string, weeks int) (int, error){    
    repos,err:=reposByOrg(client,org)
    if err != nil {
        return -1, err
    }
    
    totalCommits := 0
    for _,repo := range repos{
        if *repo.StargazersCount == 0 || *repo.Size==0 || *repo.Fork{
            continue
        }
        parts, resp, err := client.Repositories.ListParticipation(org, *repo.Name)
        if err != nil {
            return -1, err
        }
        for resp.StatusCode==202{
            time.Sleep(1000*1000*1000*1) //1 seconds
            parts, resp, err = client.Repositories.ListParticipation(org, *repo.Name)
                if err != nil {
                    return -1, err
                }
        }
        if len(parts.All)==0{
            continue
        }
        selected := parts.All[52-weeks:52]
        
        
        for _,v := range selected{
            totalCommits+=v
        }   
    }
    return totalCommits, nil   
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

func orgMembersCount(client *github.Client, org string)(int, error){
    opt:=&github.ListMembersOptions{
        ListOptions:github.ListOptions{PerPage:100},
    }
    var totalMembers int
    for {
        members, resp, err := client.Organizations.ListMembers(org, opt)
        if err!=nil{
            return 0, err
        }
        totalMembers+=len(members)
        if resp.NextPage == 0 {
            break
        }
        opt.Page = resp.NextPage
    }
    return totalMembers, nil
}