# mDNS

[![CI](https://github.com/MarlikAlmighty/mdns/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/MarlikAlmighty/mdns/actions/workflows/ci.yml) &nbsp;
[![Release to Docker Hub](https://github.com/MarlikAlmighty/mdns/actions/workflows/cd.yml/badge.svg?branch=master)](https://github.com/MarlikAlmighty/mdns/actions/workflows/cd.yml) &nbsp;
[![License](https://img.shields.io/badge/License-MIT%201.0-orange.svg)](https://github.com/MarlikAlmighty/mdns/blob/master/LICENSE) &nbsp;


mDNS is a lightweight application capable of processing DNS queries and managing DNS zones through a REST API. It is designed to provide fast and efficient operation, as well as ease of use.

With the REST API, you can add, modify, and delete DNS zones, records, and server settings. All operations are performed through HTTP requests, making DNS server management convenient and flexible.

mDNS supports various types of DNS records, such as A, CNAME, MX, TXT, and others. You can easily add and modify these records through the REST API to configure your DNS infrastructure according to your needs.

## Usage

First, get the script and make it executable:
```sh
curl -O https://raw.githubusercontent.com/MarlikAlmighty/mdns/master/install.sh
chmod +x install.sh
```

Then run it:

```sh
./install.sh
``` 

### Request examples

```sh
# Add domain
curl -X POST http://127.0.0.1:8081/dns -H 'Content-Type: application/json' \
    -d '{"domain":"example.com.", "ipv4s":["127.0.0.1"]}'

# List all domains
curl http://127.0.0.1:8081/dns

# List one domain
curl http://127.0.0.1:8081/dns/example.com.

# Update domain
curl -X PUT http://127.0.0.1:8081/dns -H 'Content-Type: application/json' \ 
    -d '{"domain":"example.com.", "ipv4s":["127.0.0.2", "127.0.0.3"]}'

# Delete domain
curl -X DELETE http://127.0.0.1:8081/dns -H 'Content-Type: application/json' \ 
    -d '{"domain":"example.com."}'
```

### API Documentation

```sh
$ swagger serve ./swagger-api/swagger.yml
```

### How to generate server

 Be careful, core methods will be overwritten.
```sh
$ swagger generate server --spec ./swagger-api/swagger.yml \ 
--target ./internal/gen -C ./swagger-templates/default-server.yml \
--template-dir ./swagger-templates --name mdns
```
