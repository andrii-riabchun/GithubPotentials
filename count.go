package githubpotentials

import (
	"github.com/google/go-github/github"
	"time"
)

func (i instance) countContributors(owner, repo string) (int, error) {
	opt := &github.CommitsListOptions{
		Since:       i.lastUpdated,
		ListOptions: github.ListOptions{PerPage: resultsPerPage},
	}
	uc := newUniqueCounter()

	for {
		commits, resp, err := i.client.Repositories.ListCommits(owner, repo, opt)
		if err != nil {
			return 0, err
		}

		for _, commit := range commits {
			if commit.Author != nil {
				uc.Add(*commit.Author.ID)
			} else if commit.Committer != nil {
				uc.Add(*commit.Committer.ID)
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return uc.Count(), nil
}

func (i instance) countCommits(owner, repo string) (int, error) {
	opt := &github.CommitsListOptions{
		Since:       i.lastUpdated,
		ListOptions: github.ListOptions{PerPage: resultsPerPage},
	}

	totalCommits := 0
	for {
		commits, resp, err := i.client.Repositories.ListCommits(owner, repo, opt)
		if err != nil {
			return 0, err
		}

		totalCommits += len(commits)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return totalCommits, nil
}

func (i instance) countStars(owner, repo string) (int, error) {
	opt := &github.ListOptions{PerPage: resultsPerPage}

	totalStars := 0

	for {
		stargazers, resp, err := i.client.Activity.ListStargazers(owner, repo, opt)
		if err != nil {
			return 0, err
		}

		filtered := filter(stargazers, isStarredAfter(i.lastUpdated))
		totalStars += len(filtered)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return totalStars, nil
}

func filter(src []*github.Stargazer, f func(*github.Stargazer) bool) []*github.Stargazer {
	var dest []*github.Stargazer
	for _, v := range src {
		if f(v) {
			dest = append(dest, v)
		}
	}
	return dest
}

func isStarredAfter(t time.Time) func(*github.Stargazer) bool {
	return func(s *github.Stargazer) bool {
		return s.StarredAt.Time.After(t)
	}
}
