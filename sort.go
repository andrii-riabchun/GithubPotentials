package githubpotentials

import (
	"github.com/artisresistance/githubpotentials/github"
	s "sort"
)

// SortCriteria enum defines all available criteria of generating result.
type SortCriteria int

func (sc SortCriteria) String() string {
	switch sc {
	case CommitsCriteria:
		return "commits"
	case StarsCriteria:
		return "stars"
	case ContributorsCriteria:
		return "contributors"
	case CombinedCriteria:
		return "sum"
	case NoCriteria:
		return "unsorted"
	default:
		return ""
	}
}

// All available sort criteria.
const (
	CommitsCriteria SortCriteria = iota
	StarsCriteria
	ContributorsCriteria
	CombinedCriteria
	NoCriteria
)

// descending sort
func sort(values []github.Repository, criteria SortCriteria) {
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

type commitsSort []github.Repository

func (s commitsSort) Len() int           { return len(s) }
func (s commitsSort) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s commitsSort) Less(i, j int) bool { return s[i].Commits > s[j].Commits }

type starsSort []github.Repository

func (s starsSort) Len() int           { return len(s) }
func (s starsSort) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s starsSort) Less(i, j int) bool { return s[i].Stars > s[j].Stars }

type contribsSort []github.Repository

func (s contribsSort) Len() int           { return len(s) }
func (s contribsSort) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s contribsSort) Less(i, j int) bool { return s[i].Contribs > s[j].Contribs }

type combinedSort []github.Repository

func (s combinedSort) Len() int      { return len(s) }
func (s combinedSort) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s combinedSort) Less(i, j int) bool {
	iStat := s[i].Commits + s[i].Contribs + s[i].Stars
	jStat := s[j].Commits + s[j].Contribs + s[j].Stars
	return iStat > jStat
}
