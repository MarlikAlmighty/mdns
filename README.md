# mDNS

### A custom dns server.

***

[![CI](https://github.com/MarlikAlmighty/mdns/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/MarlikAlmighty/mdns/actions/workflows/ci.yml) &nbsp;
[![Release to Docker Hub](https://github.com/MarlikAlmighty/mdns/actions/workflows/cd.yml/badge.svg?branch=master)](https://github.com/MarlikAlmighty/mdns/actions/workflows/cd.yml) &nbsp;
[![License](https://img.shields.io/badge/License-MIT%201.0-orange.svg)](https://github.com/MarlikAlmighty/mdns/blob/master/LICENSE) &nbsp; 

***

### Run
```sh
$ export HTTP_HOST="127.0.0.1"
$ export HTTP_PORT="8081"
$ export DNS_HOST="0.0.0.0"
$ export DNS_TCP_PORT="53"
$ export DNS_UDP_PORT="53"
$ export NAME_SERVERS="1.1.1.1,1.0.0.1,8.8.8.8,8.8.4.4"
```

### Docker
```sh
$ docker build -t marlikalmighty/mdns .
```

### Documentation: 
```sh
$ swagger serve ./swagger-api/swagger.yml
```

### How to generate server:
 Be careful, core methods will be overwritten.
```sh
$ swagger generate server --spec ./swagger-api/swagger.yml \ 
--target ./internal/gen -C ./swagger-templates/default-server.yml \
--template-dir ./swagger-templates --name mdns
```