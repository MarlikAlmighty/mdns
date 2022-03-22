# mDNS

### A custom dns server.

***

[![CI](https://github.com/MarlikAlmighty/mdns/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/MarlikAlmighty/mdns/actions/workflows/ci.yml) &nbsp;
[![Release to Docker Hub](https://github.com/MarlikAlmighty/mdns/actions/workflows/cd.yml/badge.svg?branch=master)](https://github.com/MarlikAlmighty/mdns/actions/workflows/cd.yml) &nbsp;
[![License](https://img.shields.io/badge/License-MIT%201.0-orange.svg)](https://github.com/MarlikAlmighty/mdns/blob/master/LICENSE) &nbsp; 

***

### Run
```sh
$ export REDIS_URL="redis://localhost:6379"
$ export HTTP_PORT="8081"
$ export UDP_PORT="3553"
$ export ACME_URL="https://acme-staging-v02.api.letsencrypt.org/directory"
$ export DOMAIN="example.com"
$ export IPV4="0.0.0.0"
```

### Docker
```sh
$ docker build -t mdns .
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