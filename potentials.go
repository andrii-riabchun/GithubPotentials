package githubpotentials

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const formatableQuery = "stars:>10 size:>1 pushed:>%04d-%02d-%02d"
const resultsPerPage = 100

var errAPIRateExceded = errors.New("api rate exceeded")

// ErrorHandler callback allows you to catch/log/report errors.
type ErrorHandler func(error)

// Potentials is main worker of package.
type Potentials interface {
	SearchIterator(int, ErrorHandler) RepositoryChannel
	CountStats(RepositoryChannel, ErrorHandler) RepositoryChannel
	GetAPIRates() (string, error)
}

type instance struct {
	client      *github.Client
	lastUpdated time.Time
}

// New returns new instance of package.
// token - github api token.
func New(token string, lastUpdate time.Time) Potentials {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tokenClient := oauth2.NewClient(oauth2.NoContext, tokenSource)

	return instance{
		client:      github.NewClient(tokenClient),
		lastUpdated: lastUpdate,
	}
}

func (i instance) GetAPIRates() (string, error) {
	r, _, err := i.client.RateLimits()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d of %d. Reset: %v", r.Core.Remaining, r.Core.Limit, r.Core.Reset), nil
}

func (i instance) getCoreRemainingRate() int {
	r, _, err := i.client.RateLimits()
	if err != nil {
		return -1
	}
	return r.Core.Remaining
}

// SearchIterator returns iterable channel of all search results.
// Return all repositories that were updated after specified date.
func (i instance) SearchIterator(pagesCount int, onError ErrorHandler) RepositoryChannel {
	out := make(chan RepositoryMessage)

	go func() {
		opt := &github.SearchOptions{
			Sort:        "stars",
			Order:       "asc",
			ListOptions: github.ListOptions{PerPage: resultsPerPage},
		}

		query := fmt.Sprintf(formatableQuery,
			i.lastUpdated.Year(),
			i.lastUpdated.Month(),
			i.lastUpdated.Day())
		in:=0
		for {
			result, resp, err := i.client.Search.Repositories(query, opt)
			if err != nil {
				go onError(err)
			}
			if i.getCoreRemainingRate() < resultsPerPage {
				break
			}

			for _, repo := range result.Repositories {
				casted := castRepository(repo)
				in++
				println(in, casted.Owner, casted.Name)
				out <- RepositoryMessage{&casted, resp.Remaining, nil}
			}

			if resp.NextPage == 0 || opt.Page == pagesCount-1 {
				break
			}
			opt.Page = resp.NextPage
		}
		close(out)
	}()

	return out
}

func (i instance) CountStats(in RepositoryChannel, onError ErrorHandler) RepositoryChannel {
	out := make(chan RepositoryMessage)
	go func() {
		for repo := range in {
			if repo.apiCallsRemained == 0 {
				break
			}

			joiner := new(sync.WaitGroup)
			joiner.Add(3)

			go func() {
				defer joiner.Done()
				commitsCount, err := i.countCommits(repo.repository.Owner, repo.repository.Name)
				if err != nil {
					go onError(err)
					repo.err = err
				} else {
					repo.repository.Commits = commitsCount
				}
			}()

			go func() {
				defer joiner.Done()
				starsCount, err := i.countStars(repo.repository.Owner, repo.repository.Name)
				if err != nil {
					go onError(err)
					repo.err = err
				} else {
					repo.repository.Stars = starsCount
				}
			}()

			go func() {
				defer joiner.Done()
				contribsCount, err := i.countContributors(repo.repository.Owner, repo.repository.Name)
				if err != nil {
					go onError(err)
					repo.err = err
				} else {
					repo.repository.Contribs = contribsCount
				}
			}()

			joiner.Wait()
			out <- repo
		}
		close(out)
	}()
	return out
}
