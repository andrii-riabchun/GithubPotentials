package githubpotentials

import (
	"time"
	"github.com/google/go-github/github"
)

func (i instance) countContributors(owner, repo string, date time.Time) (int, error) {
	opt := &github.CommitsListOptions{
		Since:       date,
		ListOptions: github.ListOptions{PerPage: i.resultsPerPage},
	}
	counter := newUniqueCounter()

	for {
		commits, resp, err := i.client.Repositories.ListCommits(owner, repo, opt)
		if err != nil {
			return 0, err
		}

		for _, commit := range commits {
			if commit.Author != nil {
				counter.Add(*commit.Author.ID)
			} else if commit.Committer != nil {
				counter.Add(*commit.Committer.ID)
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return counter.Count(), nil
}

func (i instance) countCommits(owner, repo string, date time.Time) (int, error) {
	opt := &github.CommitsListOptions{
		Since:       date,
		ListOptions: github.ListOptions{PerPage: i.resultsPerPage},
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

func (i instance) countStars(owner, repo string, date time.Time) (int, error) {
	opt := &github.ListOptions{PerPage: i.resultsPerPage}

	totalStars := 0

	for {
		stargazers, resp, err := i.client.Activity.ListStargazers(owner, repo, opt)
		if err != nil {
			return 0, err
		}

		filtered := filter(stargazers, filterPredicate(date))
		totalStars += len(filtered)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return totalStars, nil
}

// what if I put these inside countStars??
func filter(src []*github.Stargazer, f func(*github.Stargazer) bool) []*github.Stargazer {
	var dest []*github.Stargazer
	for _, v := range src {
		if f(v) {
			dest = append(dest, v)
		}
	}
	return dest
}

// yeah, currying!
func filterPredicate(t time.Time) func(*github.Stargazer) bool {
	return func(s *github.Stargazer) bool {
		return s.StarredAt.Time.After(t)
	}
}