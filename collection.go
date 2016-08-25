package githubpotentials

type RepositoryCollection []Repository

func (c RepositoryCollection) Trim(count int) RepositoryCollection {
	bound := count
	if len(c) < count {
		bound = len(c)
	}
	return c[:bound]
}

func (c RepositoryCollection) Sort(criteria SortCriteria) RepositoryCollection {
	sort(c, criteria)
	return c
}
