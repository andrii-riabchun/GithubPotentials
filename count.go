package githubpotentials

import (
	"github.com/artisresistance/githubpotentials/github"
)

func (i instance) countContributors(owner, repo string) (int, error) {
	uc := newUniqueCounter()
	i.client.ListCommits(owner, repo, i.lastUpdated, func(commits []github.Commit) {
		for _, commit := range commits {
			uc.Add(commit.CommitterID)
		}
	})
	return uc.Count(), nil
}

func (i instance) countCommits(owner, repo string) (int, error) {
	totalCommits := 0
	i.client.ListCommits(owner, repo, i.lastUpdated, func(commits []github.Commit) {
		totalCommits += len(commits)
	})
	return totalCommits, nil
}

func (i instance) countStars(owner, repo string) (int, error) {
	totalStars := 0
	i.client.ListStargazers(owner, repo, i.lastUpdated, func(sgs []github.Stargazer) {
		totalStars += len(sgs)
	})
	return totalStars, nil
}
