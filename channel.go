package githubpotentials

import (
	"fmt"
	"github.com/artisresistance/githubpotentials/github"
	"sync"
)

type RepositoryChannel chan *github.Repository

// Search returns iterable channel of all search results.
// Return all repositories that were updated after specified date.
func (i instance) Search(pagesCount int) RepositoryChannel {
	out := make(RepositoryChannel)

	go func() {
		query := fmt.Sprintf(formatableQuery,
			i.lastUpdated.Year(),
			i.lastUpdated.Month(),
			i.lastUpdated.Day())

		i.client.SearchRepositories(query, pagesCount, func(repos []github.Repository) {
			for _, repo := range repos {
				out <- &repo
			}
		})

		close(out)
	}()

	return out
}

func (i instance) CountStats(in RepositoryChannel) RepositoryChannel {
	out := make(RepositoryChannel)
	go func() {
		for repo := range in {
			joiner := new(sync.WaitGroup)
			joiner.Add(2)

			go func() {
				defer joiner.Done()
				commitsCount, contribsCount, err := i.countCommitsAndContributors(repo.Owner, repo.Name)
				if err != nil {
					i.log.Println(err)
					return
				}
				repo.Commits = commitsCount
				repo.Contribs = contribsCount
			}()

			go func() {
				defer joiner.Done()
				starsCount, err := i.countStars(repo.Owner, repo.Name)
				if err != nil {
					i.log.Println(err)
					return
				}
				repo.Stars = starsCount
			}()

			joiner.Wait()
			out <- repo
		}
		close(out)
	}()
	return out
}

func (in RepositoryChannel) FilterZeroStats(criteria SortCriteria) RepositoryChannel {
	out := make(RepositoryChannel)
	var isAcceptable func(repo *github.Repository) bool
	switch criteria {
	case CommitsCriteria:
		isAcceptable = func(repo *github.Repository) bool {
			return repo.Commits > 1
		}
		break
	case StarsCriteria:
		isAcceptable = func(repo *github.Repository) bool {
			return repo.Stars > 0
		}
		break
	case ContributorsCriteria:
		isAcceptable = func(repo *github.Repository) bool {
			return repo.Contribs > 0
		}
		break
	case CombinedCriteria:
		isAcceptable = func(repo *github.Repository) bool {
			return repo.Contribs+repo.Commits+repo.Stars > 1
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
		for repo := range in {
			for i := range out {
				out[i] <- repo
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
	for repo := range in {
		result = append(result, *repo)
	}
	return result
}
