package github

import "github.com/google/go-github/github"
import "context"
import "golang.org/x/oauth2"
import "log"
import "time"
import "errors"

var ResultsPerPage = 100

var errAPIRateExceded = errors.New("api rate exceeded")

type Client interface {
	SearchRepositories(query string, pages int, onFetch func([]Repository))
	ListCommits(owner, repo string, since time.Time, onFetch func([]Commit))
	ListStargazers(owner, repo string, since time.Time, onFetch func([]Stargazer))
	APIRates() (remaining int, reset time.Time, err error)
}

func NewClient(token string, log *log.Logger) Client {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	return client{
		log,
		github.NewClient(oauthClient),
	}
}

type client struct {
	log    *log.Logger
	client *github.Client
}

func (c client) APIRates() (int, time.Time, error) {
	rates, _, err := c.client.RateLimits()
	return rates.Core.Remaining, rates.Core.Reset.Time, err
}

func (c client) SearchRepositories(query string, pages int, onFetch func([]Repository)) {
	opt := &github.SearchOptions{
		Sort:        "stars",
		Order:       "asc",
		ListOptions: github.ListOptions{PerPage: ResultsPerPage},
	}
	repos := make([]Repository, ResultsPerPage)
	for {
		result, resp, err := c.client.Search.Repositories(query, opt)
		if err != nil {
			c.log.Println(err)
			break
		}

		for i, repo := range result.Repositories {
			repos[i] = castRepository(repo)
		}

		onFetch(repos[:len(result.Repositories)])

		if resp.NextPage == 0 || opt.Page == pages-1 {
			break
		}
		opt.Page = resp.NextPage
	}
}

func (c client) ListCommits(owner, repo string, since time.Time, onFetch func([]Commit)) {
	opt := &github.CommitsListOptions{
		Since:       since,
		ListOptions: github.ListOptions{PerPage: ResultsPerPage},
	}
	//TODO move out of loop
	out := make([]Commit, ResultsPerPage)
	for {
		commits, resp, err := c.client.Repositories.ListCommits(owner, repo, opt)
		if err != nil {
			c.log.Println(err)
			break
		}

		for i, commit := range commits {
			if commit.Author != nil {
				out[i] = Commit{*commit.Author.ID}
			} else if commit.Committer != nil {
				out[i] = Commit{*commit.Committer.ID}
			}
		}

		onFetch(out[:len(commits)])

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
}

func (c client) ListStargazers(owner, repo string, since time.Time, onFetch func([]Stargazer)) {
	opt := &github.ListOptions{PerPage: ResultsPerPage}

	out := make([]Stargazer, ResultsPerPage)
	for {
		stargazers, resp, err := c.client.Activity.ListStargazers(owner, repo, opt)
		if err != nil {
			c.log.Println(err)
			break
		}

		total := 0
		for _, sg := range stargazers {
			if time := sg.StarredAt.Time; time.After(since) {
				out[total] = Stargazer{time}
				total++
			}
		}

		onFetch(out[:total])

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
}
