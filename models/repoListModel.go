package models

type RepoList struct {
    Criteria    string      `json:"criteria,omitempty"`
    Items       []RepoEntry `json:"items,omitempty"`
}

type RepoEntry struct{
    FullName string         `json:"fullname,omitempty"`
    SortingCriteria int     `json:"value"`
}

type RepoSort []RepoEntry

func (slice RepoSort) Len() int {
    return len(slice)
}

func (slice RepoSort) Less(i, j int) bool {
    return slice[i].SortingCriteria > slice[j].SortingCriteria;
}

func (slice RepoSort) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}

