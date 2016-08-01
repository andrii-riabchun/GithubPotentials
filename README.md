# Github Potentials

Github Potentials command line tool is package that allows you to find rising Github repositories in 3 steps:
  - Fetch 1k recently updated repositories.
  - Count stats for last *n* hours/days: new stars, unique contributors and commits.
  - Sort by selected criteria and take the best.

### Command line tool
Under [cmd](https://github.com/ArtIsResistance/GithubPotentials/tree/master/ghp/main.go) directory you can see an example of using this package.

`go get github.com/artisresistance/githubpotentials/ghp`

You must provide config file that contains your Github API secret token.
Example output you can find [here](https://githubpotentials.azure.net/data.json).