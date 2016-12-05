package githubpotentials

import (
	"github.com/artisresistance/githubpotentials/github"
)

func (i instance) countCommitsAndContributors(owner, repo string) (int, int, error) {
	totalCommits := 0
	contributorsUC := newUniqueCounter()
	i.client.ListCommits(owner, repo, i.lastUpdated, func(commits []github.Commit) {
		totalCommits += len(commits)
		for _, commit := range commits {
			contributorsUC.Add(commit.CommitterID)
		}
	})
	return totalCommits, contributorsUC.Count(), nil
}

func (i instance) countStars(owner, repo string) (int, error) {
	totalStars := 0
	i.client.ListStargazers(owner, repo, i.lastUpdated, func(sgs []github.Stargazer) {
		totalStars += len(sgs)
	})
	return totalStars, nil
}
