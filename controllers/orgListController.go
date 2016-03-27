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

// GetOrgList GET /orgs/:criteria/:weeks
func GetOrgList(params martini.Params, req *http.Request) (int, []byte){
    weeks, err := strconv.Atoi(params["weeks"])
    if err != nil || weeks<1 || weeks>52 {
        return http.StatusBadRequest, nil
    }
    
    date := time.Now().AddDate(0,0,-weeks*7)      
    query:="repos:>1 type:org created:>"+helpers.FormatDate(date)
 
    result := models.OrgList{}
    switch params["criteria"] {
    case "total":
        orgs, err := listOrgByTotalCommits(query,weeks)
        if err != nil {
            return http.StatusInternalServerError, nil
        }
        result.Criteria="total"
        result.Items=orgs
    case "avg":
        orgs, err := listOrgByAvgCommits(query,weeks)
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
    cache.Add(req.RequestURI, resp, expiration())
    return http.StatusOK, resp
}



func listOrgByTotalCommits(query string, weeks int) ([]models.OrgEntry, error){
    var orgResult models.OrgSort
    orgs,err := searchOrgs(query)
    if err != nil {
        return nil, err
    }
    for _,org:= range orgs{                   
        orgCommits, err:=getCommitsByOrg(*org.Login,weeks)
        if err != nil {
            return nil,err
        }
        if orgCommits==0 {
            continue
        }           
        orgCommitsFl := float32(orgCommits)
        result:=models.OrgEntry{
            Name: *org.Login,
            SortingCriteria: orgCommitsFl,
        }
        orgResult = append(orgResult, result)
        
    }
    sort.Sort(orgResult)
    return orgResult, nil
}

func listOrgByAvgCommits(query string, weeks int) ([]models.OrgEntry, error){
    var orgResult models.OrgSort
    orgs,err := searchOrgs(query)
    if err != nil {
        return nil,err
    }
    for _,org:= range orgs{              
        orgCommits, err:=getCommitsByOrg(*org.Login,weeks)
        if err != nil {
            return nil, err
        }
        if orgCommits==0 {
            continue
        }
        
        members, err := orgMembersCount(*org.Login)
        if err != nil {
            return nil, err
        }
        if members==0{
            continue
        }
        avgCommits := float32(orgCommits)/float32(members)
        
        result:=models.OrgEntry{
            Name: *org.Login,
            SortingCriteria: avgCommits,
        }
        
        orgResult = append(orgResult, result)
    }
    sort.Sort(orgResult)
    return orgResult, nil
}

func searchOrgs(query string) ([]github.Organization,error){
    if cached,is:=cache.Get("orgsearch:"+query);is{
        return cached.([]github.Organization), nil
    }
    i:=0
    opt:=&github.SearchOptions{
        Sort:   "followers",
        Order:  "asc",
        ListOptions: github.ListOptions{
            PerPage:100,
        },
    }
    var organizations []github.Organization
    for {           
        i++
        orgs, resp, err := client.Search.Users(query, opt)
        if err!=nil{
            return nil, err
        }
        for _,orgUser:= range orgs.Users{
            org, _, err := client.Organizations.Get(*orgUser.Login)
            if err != nil{
                continue
            } 
            organizations = append(organizations, *org)            
        }
        if resp.NextPage == 0 || i==orgPagesCount {
            break
        }
        opt.Page = resp.NextPage
    }
    cache.Add("orgsearch:"+query, organizations, expiration())
    return organizations, nil
}


func getCommitsByOrg(org string, weeks int) (int, error){ 
    var repos []github.Repository    
    repos,err:=reposByUser(org)
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


func orgMembersCount(org string)(int, error){
    if cached,is:=cache.Get("org-members:"+org);is{
        return cached.(int), nil
    }
    
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
    cache.Add("org-members:"+org, totalMembers, expiration())
    return totalMembers, nil
}