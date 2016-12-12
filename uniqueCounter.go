package githubpotentials

type uniqueCounter struct {
	items map[interface{}]struct{}
}

func newUniqueCounter() uniqueCounter {
	return uniqueCounter{
		items: make(map[interface{}]struct{}),
	}
}

func (uc uniqueCounter) Add(item interface{}) {
	uc.items[item] = struct{}{}
}

func (uc uniqueCounter) Count() int {
	return len(uc.items)
}
