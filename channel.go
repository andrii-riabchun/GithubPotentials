package githubpotentials

import (
    "sync"
    "fmt"
    "github.com/google/go-github/github"
)

type RepositoryMessage struct {
	repository       *Repository
	apiCallsRemained int
	err              error
}

type RepositoryChannel chan RepositoryMessage

// Search returns iterable channel of all search results.
// Return all repositories that were updated after specified date.
func (i instance) Search(pagesCount int, onError ErrorHandler) RepositoryChannel {
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
			return repoMsg.repository.Contribs > 1
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

func (in RepositoryChannel) Split(count int) []RepositoryChannel{
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

func (in RepositoryChannel) Dump(onError ErrorHandler) RepositoryCollection {
	var result []Repository
	for repoMsg := range in {
		if repoMsg.err != nil {
			onError(repoMsg.err)
			continue
		}

		result = append(result, *repoMsg.repository)

		if repoMsg.apiCallsRemained == 0 {
			onError(errAPIRateExceded)
			break
		}
	}

	return result
}