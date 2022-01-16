# mDns

### A custom dns server.

***

[![CI](https://github.com/MarlikAlmighty/mDns/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/MarlikAlmighty/mDns/actions/workflows/ci.yml) &nbsp;
[![Release to Docker Hub](https://github.com/MarlikAlmighty/mDns/actions/workflows/cd.yml/badge.svg?branch=master)](https://github.com/MarlikAlmighty/mDns/actions/workflows/cd.yml) &nbsp;
[![License](https://img.shields.io/badge/License-MIT%201.0-orange.svg)](https://github.com/MarlikAlmighty/mDns/blob/master/LICENSE) &nbsp; 

***

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
--template-dir ./swagger-templates --name mdns
```