# mDNS

[![CI](https://github.com/MarlikAlmighty/mdns/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/MarlikAlmighty/mdns/actions/workflows/ci.yml) &nbsp;
[![Release to Docker Hub](https://github.com/MarlikAlmighty/mdns/actions/workflows/cd.yml/badge.svg?branch=master)](https://github.com/MarlikAlmighty/mdns/actions/workflows/cd.yml) &nbsp;
[![License](https://img.shields.io/badge/License-MIT%201.0-orange.svg)](https://github.com/MarlikAlmighty/mdns/blob/master/LICENSE) &nbsp;


**mDNS** is a lightweight application capable of processing dns queries and managing dns zones through a rest api. With the rest api, you can add, modify, and delete dns zones, records, and server settings. All operations are performed through http requests, making dns server management convenient and flexible.

## Usage

Open 53 ports in your firewall:
```sh
sudo ufw enable
sudo ufw allow 53/tcp
sudo ufw allow 53/udp
```

Get the script:
```sh
curl -O https://raw.githubusercontent.com/MarlikAlmighty/mdns/master/ubuntu-server-install.sh
```

Make it executable:
```sh
sudo chmod +x ubuntu-server-install.sh
```

Then run it:

```sh
sudo ./ubuntu-server-install.sh
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
