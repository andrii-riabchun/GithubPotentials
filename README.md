# Github Potentials

Github Potentials is service that allows you to:
  - Show stats for selected repository.
  - Browse the most popular of the the newest repositories.
  - List perspective organizations with great leap in recent activity.

It can take a long time to create list of repos and organizations - up to five minutes - so please be patient, sit back and relax while it processes your request.

### Architecture

Frontend is built as Single Page Application with JS, while Golang-based backend serves some kind of REST Web API.

* Foundation - beautiful page layout framework.
* jQuery - manipulate browser DOM-tree with JS.
* Chart.js - chart drawing.
* go-martini - fast and simple golang web framowork.
* facebookgo/inmem - who needs memcached?