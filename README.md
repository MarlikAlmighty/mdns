# Library

### A simple example of a server API, with a clean architecture.

***

[![CI](https://github.com/MarlikAlmighty/library/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/MarlikAlmighty/library/actions/workflows/ci.yml) &nbsp;
[![Release to Docker Hub](https://github.com/MarlikAlmighty/library/actions/workflows/cd.yml/badge.svg?branch=master)](https://github.com/MarlikAlmighty/library/actions/workflows/cd.yml) &nbsp;
[![License](https://img.shields.io/badge/License-MIT%201.0-orange.svg)](https://github.com/MarlikAlmighty/library/blob/master/LICENSE) &nbsp; 

***

I chose the library as a reference. We can manage books, bookcases.

***

### How to run
```sh
$ docker stack deploy postgres -c postgresql.yml

$ export PREFIX=LIBRARY

$ export LIBRARY_HTTP_PORT=3000

$ export LIBRARY_DBCONNECT=postgres://user:password@localhost:5432/library?sslmode=disable

$ export LIBRARY_MIGRATE=true

$ export LIBRARY_PATH_TO_MIGRATE="migrations"

$ go run ./cmd/...
```

### Docker
```sh
$ docker build -t library .
```

### Documentation: 
```sh
$ swagger serve ./swagger-api/swagger.yml
```

### How to generate server:
```sh
$ swagger generate server --spec ./swagger-api/swagger.yml \ 
--target ./internal/gen -C ./swagger-templates/default-server.yml \
--template-dir ./swagger-templates --name library
```