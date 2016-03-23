package helpers

import (
    "reflect"
    "net/url"
    "fmt"
    "github.com/google/go-github/github"
    "github.com/google/go-querystring/query"
)

type Stargazer struct {
	StarredAt  *github.Timestamp  `json:"starred_at,omitempty"`
	User *github.User `json:"user,omitempty"`
}
// ListStargazers with 'starred_at' field
func ListStargazers(s *github.Client, owner, repo string, opt *github.ListOptions) ([]Stargazer, *github.Response, error) {
	u := fmt.Sprintf("repos/%s/%s/stargazers", owner, repo)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

    req.Header.Set("Accept","application/vnd.github.v3.star+json")
	stargazers := new([]Stargazer)
	resp, err := s.Do(req, stargazers)
	if err != nil {
		return nil, resp, err
	}

	return *stargazers, resp, err
}

func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}