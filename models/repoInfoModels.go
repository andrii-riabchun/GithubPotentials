package models

type RepoStars struct{
    Owner       string      `json:"owner,omitempty"`
    Repo        string      `json:"repo,omitempty"`
    DayCount    int         `json:"days,omitempty"`
    TotalStars  int         `json:"total,omitempty"`
    DayStars    []int       `json:"stars,omitempty"`   
}

type RepoCommits struct{
    Owner           string      `json:"owner,omitempty"`
    Repo            string      `json:"repo,omitempty"`
    DayCount        int         `json:"days,omitempty"`
    TotalCommits    int         `json:"total,omitempty"`
    DayCommits      []int       `json:"commits,omitempty"` 
}

type RepoContributors struct{
    Owner           string      `json:"owner,omitempty"`
    Repo            string      `json:"repo,omitempty"`
    DayCount        int         `json:"days,omitempty"`
    UniqueContribs  int         `json:"total,omitempty"`
    DayContribs     []int       `json:"contribs,omitempty"` 
}