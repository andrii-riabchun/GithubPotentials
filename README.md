# Github Potentials

Github Potentials is package that allows you to get rising stars of Github in three steps:
  - Search repositories that were recently updated.
  - Count stats for each of them - new stars, unique contributors and commits.
  - Sort by selected criteria and take top N repositories.

It can take a long time to create list of repos - up to five minutes - so please be patient, sit back and relax while it executes.

### Architecture

Frontend is built as Single Page Application with JS, while Golang-based backend serves some kind of REST Web API.

* Foundation - beautiful page layout framework.
* jQuery - manipulate browser DOM-tree with JS.
* Chart.js - chart drawing.
* go-martini - fast and simple golang web framowork.
* facebookgo/inmem - who needs memcached?