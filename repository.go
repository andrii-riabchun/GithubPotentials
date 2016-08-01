package githubpotentials

import (
	"github.com/google/go-github/github"
)

type RepositoryChannel <-chan RepositoryMessage

func (in RepositoryChannel) FilterZeroStats(criteria SortCriteria) RepositoryChannel {
	out := make(chan RepositoryMessage)

	isAcceptable := func(repoMsg RepositoryMessage) bool { return false }
	switch criteria {
	case CommitsCriteria:
		isAcceptable = func(repoMsg RepositoryMessage) bool {
			return repoMsg.repository.Commits > 0
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

type RepositoryCollection []Repository

func (c RepositoryCollection) Sort(criteria SortCriteria) RepositoryCollection {
	sort(c, criteria)
	return c
}

type RepositoryMessage struct {
	repository       *Repository
	apiCallsRemained int
	err              error
}

type Repository struct {
	Owner       string
	Name        string
	Description string
	Homepage    string
	License     string
	Language    string
	Commits     int
	Stars       int
	Contribs    int
}

func castRepository(src github.Repository) Repository {
	result := Repository{
		Owner: *src.Owner.Login,
		Name:  *src.Name,
	}
	if src.Description != nil {
		result.Description = *src.Description
	}
	if src.Homepage != nil {
		result.Homepage = *src.Homepage
	}
	if src.Language != nil {
		result.Language = *src.Language
	}
	if src.License != nil && src.License.Name != nil {
		result.License = *src.License.Name
	}

	return result
}
