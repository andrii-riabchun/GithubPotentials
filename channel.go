package githubpotentials

import (
	"fmt"
	"github.com/artisresistance/githubpotentials/github"
	"sync"
)

type RepositoryMessage struct {
	repository       *github.Repository
	apiCallsRemained int
	err              error
}

type RepositoryChannel chan RepositoryMessage

// Search returns iterable channel of all search results.
// Return all repositories that were updated after specified date.
func (i instance) Search(pagesCount int) RepositoryChannel {
	out := make(chan RepositoryMessage)

	go func() {
		query := fmt.Sprintf(formatableQuery,
			i.lastUpdated.Year(),
			i.lastUpdated.Month(),
			i.lastUpdated.Day())

		i.client.SearchRepositories(query, pagesCount, func(repos []github.Repository) {
			for _, repo := range repos {
				//TODO remove RepositoryMessage as type
				out <- RepositoryMessage{&repo, 1000, nil}
			}
		})

		close(out)
	}()

	return out
}

func (i instance) CountStats(in RepositoryChannel) RepositoryChannel {
	out := make(chan RepositoryMessage)
	go func() {
		for repo := range in {
			if repo.apiCallsRemained == 0 {
				break
			}

			joiner := new(sync.WaitGroup)
			joiner.Add(2)

			go func() {
				defer joiner.Done()
				commitsCount, contribsCount, err := i.countCommitsAndContributors(repo.repository.Owner, repo.repository.Name)
				if err != nil {
					i.log.Println(err)
					repo.err = err
				} else {
					repo.repository.Commits = commitsCount
					repo.repository.Contribs = contribsCount
				}
			}()

			go func() {
				defer joiner.Done()
				starsCount, err := i.countStars(repo.repository.Owner, repo.repository.Name)
				if err != nil {
					i.log.Println(err)
					repo.err = err
				} else {
					repo.repository.Stars = starsCount
				}
			}()

			joiner.Wait()
			out <- repo
		}
		close(out)
	}()
	return out
}

func (in RepositoryChannel) FilterZeroStats(criteria SortCriteria) RepositoryChannel {
	out := make(chan RepositoryMessage)
	var isAcceptable func(repoMsg RepositoryMessage) bool
	switch criteria {
	case CommitsCriteria:
		isAcceptable = func(repoMsg RepositoryMessage) bool {
			return repoMsg.repository.Commits > 1
		}
		break
	case StarsCriteria:
		isAcceptable = func(repoMsg RepositoryMessage) bool {
			return repoMsg.repository.Stars > 0
		}
		break
	case ContributorsCriteria:
		isAcceptable = func(repoMsg RepositoryMessage) bool {
			return repoMsg.repository.Contribs > 0
		}
		break
	case CombinedCriteria:
		isAcceptable = func(repoMsg RepositoryMessage) bool {
			return repoMsg.repository.Contribs+
				repoMsg.repository.Commits+
				repoMsg.repository.Stars > 0
		}
		break
	}

	go func() {
		for repo := range in {
			if isAcceptable(repo) {
				out <- repo
			}
		}
		close(out)
	}()
	return out
}

func (in RepositoryChannel) Split(count int) []RepositoryChannel {
	out := make([]RepositoryChannel, count)
	for i := range out {
		out[i] = make(RepositoryChannel)
	}
	go func() {
		for msg := range in {
			for i := range out {
				out[i] <- msg
			}
		}
		for i := range out {
			close(out[i])
		}
	}()
	return out
}

func (in RepositoryChannel) Dump() RepositoryCollection {
	var result []github.Repository
	for repoMsg := range in {
		if repoMsg.err != nil {
			continue
		}

		result = append(result, *repoMsg.repository)

		if repoMsg.apiCallsRemained == 0 {
			break
		}
	}

	return result
}
