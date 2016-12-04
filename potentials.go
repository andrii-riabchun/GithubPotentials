package githubpotentials

import (
	"time"

	"github.com/artisresistance/githubpotentials/github"
	"log"
)

const formatableQuery = "stars:>10 size:>1 pushed:>%04d-%02d-%02d"

// ErrorHandler callback allows you to catch/log/report errors.
type ErrorHandler func(error)

// Potentials is main worker of package.
type Potentials interface {
	Search(int, ErrorHandler) RepositoryChannel
	CountStats(RepositoryChannel, ErrorHandler) RepositoryChannel
	APIRates() (remaining int, reset time.Time, err error)
}

type instance struct {
	client      github.Client
	lastUpdated time.Time
}

// New returns new instance of package.
// token - github api token.
func New(token string, since time.Time, log *log.Logger) Potentials {
	return instance{
		client:      github.NewClient(token, log),
		lastUpdated: since,
	}
}

func (i instance) APIRates() (int, time.Time, error) {
	return i.client.APIRates()
}
