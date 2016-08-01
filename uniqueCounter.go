package githubpotentials

type uniqueCounter struct {
	items map[int]struct{}
}

func newUniqueCounter() uniqueCounter {
	return uniqueCounter{items: make(map[int]struct{})}
}

func (uc uniqueCounter) Add(i int) {
	uc.items[i] = struct{}{}
}

func (uc uniqueCounter) Count() int {
	return len(uc.items)
}
