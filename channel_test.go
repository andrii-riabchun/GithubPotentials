package githubpotentials

import "testing"
import "github.com/artisresistance/githubpotentials/github"
import "time"

func TestSearch(t *testing.T) {
	ch := newTestInstance().Search(10)
	i := 0
	for repo := range ch {
		if repo == nil {
			t.Error("received repo is nil")
		}
		i++
	}
	if i != 1000 {
		t.Errorf("i==%d, expected to be %d", i, 1000)
	}
}

func TestCountStats(t *testing.T) {
	inst := newTestInstance()
	in := make(RepositoryChannel)
	go func() {
		for repo := range inst.CountStats(in) {
			if repo.Commits != itemsPerRepo {
				t.Errorf("repo.Commits==%d, expected %d", repo.Commits, itemsPerRepo)
			}
			if repo.Contribs != 15 {
				t.Errorf("repo.Contribs==%d, expected %d", repo.Contribs, 15)
			}
			if repo.Stars != itemsPerRepo {
				t.Errorf("repo.Stars==%d, expected %d", repo.Stars, itemsPerRepo)
			}
		}
	}()
	for i := 0; i < 10; i++ {
		in <- &github.Repository{}
	}
}

func TestAPIRates(t *testing.T) {
	mustReset := time.Now().Add(time.Hour)
	mustRemaining := 5000 - 3*100*5
	inst := newTestInstance()

	ch := inst.Search(5)
	inst.CountStats(ch)

	remaining, reset, err := inst.APIRates()
	if remaining == mustRemaining {
		t.Errorf("remaining==%d, expected %d", remaining, mustRemaining)
	}
	if reset == mustReset {
		t.Errorf("reset==%d, expected %d", reset, mustReset)
	}
	if err != nil {
		t.Error("err expected to be nil")
	}
}
