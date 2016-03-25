package models

type OrgList struct {
    Criteria    string      `json:"criteria,omitempty"`
    Items       []OrgEntry  `json:"items,omitempty"`
}

type OrgEntry struct{
    Name string             `json:"name,omitempty"`
    SortingCriteria float32 `json:"value"`
}

type OrgSort []OrgEntry

func (slice OrgSort) Len() int {
    return len(slice)
}

func (slice OrgSort) Less(i, j int) bool {
    return slice[i].SortingCriteria > slice[j].SortingCriteria;
}

func (slice OrgSort) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}