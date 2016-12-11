package githubpotentials

import "log"
import "github.com/artisresistance/githubpotentials/github"
import "time"
import "fmt"
import "io/ioutil"

const searchItemsPerRequest = 100
const itemsPerRepo = 50

type clientMock struct {
	begin        time.Time
	apiRemaining int
	repos        []github.Repository
}

func (m clientMock) repo(owner, name string) *github.Repository {
	for _, repo := range m.repos {
		if repo.Owner == owner && repo.Name == name {
			return &repo
		}
	}
	return nil
}

func (m clientMock) SearchRepositories(query string, pages int, onFetch func([]github.Repository)) {
	for i := 0; i < pages*searchItemsPerRequest; i += searchItemsPerRequest {
		m.apiRemaining--
		onFetch(m.repos[i : i+searchItemsPerRequest])
	}
}

func (m clientMock) ListCommits(owner, repo string, since time.Time, onFetch func([]github.Commit)) {
	commits := make([]github.Commit, itemsPerRepo)
	for i := 0; i < itemsPerRepo; i++ {
		commits[i].CommitterID = i % 15
	}
	m.apiRemaining--
	onFetch(commits)
}

func (m clientMock) ListStargazers(owner, repo string, since time.Time, onFetch func([]github.Stargazer)) {
	sgs := make([]github.Stargazer, itemsPerRepo)
	for i := 0; i < itemsPerRepo; i++ {
		sgs[i].StarredAt = time.Now()
	}
	m.apiRemaining--
	onFetch(sgs)
}

func (m clientMock) APIRates() (remaining int, reset time.Time, err error) {
	return m.apiRemaining, m.begin.Add(time.Hour), nil
}

func newTestInstance() Potentials {
	mock := clientMock{apiRemaining: 5000}
	mock.repos = make([]github.Repository, 1000)
	for i, repo := range mock.repos {
		repo.Name = fmt.Sprintf("%d", i)
		repo.Owner = fmt.Sprintf("%d", i%100)
		repo.Description = fmt.Sprintf("desc for %s/%s", repo.Owner, repo.Name)
		repo.Homepage = fmt.Sprintf("homepage for %s/%s", repo.Owner, repo.Name)
		repo.License = "MIT"
	}
	return instance{
		log:         log.New(ioutil.Discard, "", 0),
		client:      mock,
		lastUpdated: time.Now(),
	}
}
