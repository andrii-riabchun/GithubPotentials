package githubpotentials

import (
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"sync"
	"time"
)

var errAPIRateExceded = errors.New("api rate exceeded")

type Potentials interface {
	SearchIterator(date time.Time) RepositoriesChannel
	SetCriterias(in RepositoriesChannel, date time.Time) RepositoriesChannel
	GetAPIRates() string
}

type instance struct {
	client         *github.Client
	resultsPerPage int
}

// New returns new instance of package.
// token - github api token.
func New(token string, resultsPerPage int) Potentials {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tokenClient := oauth2.NewClient(oauth2.NoContext, tokenSource)

	return instance{
		client:         github.NewClient(tokenClient),
		resultsPerPage: resultsPerPage,
	}
}

func (i instance) GetAPIRates() string{
	r, _, err := i.client.RateLimits()
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("%d of %d. Reset: %v", r.Core.Remaining, r.Core.Limit, r.Core.Reset)
}

// SearchIterator returns iterable channel of all search results
func (i instance) SearchIterator(date time.Time) RepositoriesChannel {
	out := make(chan RepositoryMessage)

	go func() {
		opt := &github.SearchOptions{
			Sort:        "stars",
			Order:       "asc",
			ListOptions: github.ListOptions{PerPage: i.resultsPerPage},
		}

		query := fmt.Sprintf("stars:>10 size:>1 pushed:>%04d-%02d-%02d", date.Year(), date.Month(), date.Day())

		for {
			result, resp, err := i.client.Search.Repositories(query, opt)
			if err != nil {
				fmt.Println(err)	
			}

			for _, repo := range result.Repositories {
				casted := castRepository(repo)
				out <- RepositoryMessage{&casted, resp.Remaining, nil}
			}

			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage
		}
		close(out)
	}()

	return out
}

func (i instance) SetCriterias(in RepositoriesChannel, date time.Time) RepositoriesChannel {
	out := make(chan RepositoryMessage)
	go func() {
		for repo := range in {
			//bypass
			if repo.err != nil {
				fmt.Println(repo.err.Error())
				out <- repo
				continue
			}

			joiner := new(sync.WaitGroup)
			joiner.Add(3)

			go func() {
				defer joiner.Done()
				commitsCount, err := i.countCommits(repo.repository.owner, repo.repository.name, date)
				if err != nil {
					repo.err = err
				} else {
					repo.repository.commits = commitsCount
				}
			}()

			go func() {
				defer joiner.Done()
				starsCount, err := i.countStars(repo.repository.owner, repo.repository.name, date)
				if err != nil {
					repo.err = err
				} else {
					repo.repository.stars = starsCount
				}
			}()

			go func() {
				defer joiner.Done()
				contribsCount, err := i.countContributors(repo.repository.owner, repo.repository.name, date)
				if err != nil {
					repo.err = err
				} else {
					repo.repository.contribs = contribsCount
				}
			}()

			joiner.Wait()
			out <- repo
		}
		close(out)
	}()
	return out
}
