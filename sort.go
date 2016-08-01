package githubpotentials

import (
    s "sort"
)

// SortCriteria enum defines all available criterias of generating result.
type SortCriteria int

// All available sort criterias.
const (
	CommitsCriteria      SortCriteria = iota
	StarsCriteria        SortCriteria = iota
	ContributorsCriteria SortCriteria = iota
	CombinedCriteria     SortCriteria = iota
	NoCriteria           SortCriteria = iota
)

func sort(values []Repository, criteria SortCriteria){
    switch criteria {
    case CommitsCriteria:
        s.Sort(commitsSort(values))
        break
    case StarsCriteria:
        s.Sort(starsSort(values))
        break
    case ContributorsCriteria:
        s.Sort(contribsSort(values))
        break
    case CombinedCriteria:
        s.Sort(combinedSort(values))
        break
    }
}

type commitsSort []Repository
func (s commitsSort) Len() int { return len(s) }
func (s commitsSort) Swap(i,j int) { s[i],s[j] = s[j], s[i] }
func (s commitsSort) Less(i,j int) bool { return s[i].commits < s[j].commits }

type starsSort []Repository
func (s starsSort) Len() int { return len(s) }
func (s starsSort) Swap(i,j int) { s[i],s[j] = s[j], s[i] }
func (s starsSort) Less(i,j int) bool { return s[i].stars < s[j].stars }

type contribsSort []Repository
func (s contribsSort) Len() int { return len(s) }
func (s contribsSort) Swap(i,j int) { s[i],s[j] = s[j], s[i] }
func (s contribsSort) Less(i,j int) bool { return s[i].contribs < s[j].contribs }

type combinedSort []Repository
func (s combinedSort) Len() int { return len(s) }
func (s combinedSort) Swap(i,j int) { s[i],s[j] = s[j], s[i] }
func (s combinedSort) Less(i,j int) bool {
    iStat := s[i].commits + s[i].contribs + s[i].stars
    jStat := s[j].commits + s[j].contribs + s[j].stars
    return iStat < jStat
}
