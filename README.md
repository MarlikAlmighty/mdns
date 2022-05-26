# mDNS

### A custom dns server.

***

[![CI](https://github.com/MarlikAlmighty/mdns/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/MarlikAlmighty/mdns/actions/workflows/ci.yml) &nbsp;
[![Release to Docker Hub](https://github.com/MarlikAlmighty/mdns/actions/workflows/cd.yml/badge.svg?branch=master)](https://github.com/MarlikAlmighty/mdns/actions/workflows/cd.yml) &nbsp;
[![License](https://img.shields.io/badge/License-MIT%201.0-orange.svg)](https://github.com/MarlikAlmighty/mdns/blob/master/LICENSE) &nbsp; 

***

### Run
```sh
$ export HTTP_PORT="8081"
$ export NAME_SERVERS="1.1.1.1:53,1.0.0.1:53,8.8.8.8:53,8.8.4.4:53"
$ export DNS_HOST="127.0.0.1"
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