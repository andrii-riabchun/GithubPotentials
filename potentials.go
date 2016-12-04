package githubpotentials

import (
	"context"
	"errors"
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
	Search(int, ErrorHandler) RepositoryChannel
	CountStats(RepositoryChannel, ErrorHandler) RepositoryChannel
	GetAPIRates() (int, time.Time, error)
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
	tokenClient := oauth2.NewClient(context.Background(), tokenSource)

	return instance{
		client:      github.NewClient(tokenClient),
		lastUpdated: lastUpdate,
	}
}

func (i instance) GetAPIRates() (int, time.Time, error) {
	r, _, err := i.client.RateLimits()
	if err != nil {
		return -1, time.Time{}, err
	}
	return r.Core.Limit - r.Core.Remaining, r.Core.Reset.Time, nil
}

func (i instance) getCoreRemainingRate() int {
	r, _, err := i.client.RateLimits()
	if err != nil {
		return -1
	}
	return r.Core.Remaining
}
