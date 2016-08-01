package githubpotentials

import (
	"github.com/google/go-github/github"
)

type RepositoriesChannel <-chan RepositoryMessage

func (in RepositoriesChannel) Dump() RepositoryCollection {
	var result []Repository
	for repoMsg := range in {
		if repoMsg.err != nil {
			continue
		}

		result = append(result, *repoMsg.repository)

		println(len(result))

		if repoMsg.apiCallsRemained == 0 {
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
	owner       string
	name        string
	description string
	homepage    string
	license     string
	language    string
	commits     int
	stars       int
	contribs    int
}

func castRepository(src github.Repository) Repository {
	result := Repository{
		owner: *src.Owner.Login,
		name: *src.Name,
	}
	if src.Description 	!= nil { result.description = *src.Description }
	if src.Homepage 	!= nil { result.homepage = *src.Homepage }
	if src.Language 	!= nil { result.language = *src.Language}
	if src.License 		!= nil && src.License.Name != nil { result.license = *src.License.Name }

	return result
}
