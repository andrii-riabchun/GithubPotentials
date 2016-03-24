package models

type RepoInfo struct{
    Owner           string      `json:"owner,omitempty"`
    Repo            string      `json:"repo,omitempty"`
    DayCount        int         `json:"days,omitempty"`
    Stars           int         `json:"stars,omitempty"`
    StarsData       []int       `json:"starsdata,omitempty"`   
    Commits         int         `json:"commits,omitempty"`
    CommitsData     []int       `json:"commitsdata,omitempty"` 
    Contribs        int         `json:"contribs,omitempty"`
    ContribsData    []int       `json:"contribsdata,omitempty"` 
}