package main

import (
	"context"
	"encoding/json"
	"github.com/MarlikAlmighty/mdns/internal/app"
	"github.com/MarlikAlmighty/mdns/internal/config"
	"github.com/MarlikAlmighty/mdns/internal/data"
	"github.com/MarlikAlmighty/mdns/internal/dns"
	"github.com/MarlikAlmighty/mdns/internal/dump"
	"github.com/MarlikAlmighty/mdns/internal/gen/models"
	"github.com/MarlikAlmighty/mdns/internal/gen/restapi"
	"github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations"
	apiAdd "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/add"
	apiDelete "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/delete"

	apiShow "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/show"

	apiList "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/list"

	apiUpdate "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/update"

	"log"
	"strconv"

	"github.com/go-openapi/loads"
)

func main() {

	const dumpName = "DUMP"

	// for to catch panic
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered func, we got panic:", r)
		}
	}()

	// get the configuration for the application through ENV
	c := config.New()
	if err := c.GetEnv(); err != nil {
		log.Fatal("get environment keys: ", err)
	}

	// start store
	s, err := dump.New(c)
	if err != nil {
		log.Fatal(err)
	}

	// ping redis
	if _, err := s.Client.Ping(context.Background()).Result(); err != nil {
		log.Fatal(err)
	}

	// retrieve data from redis
	var res []byte
	if res, err = s.Pop(dumpName); err != nil {
		log.Fatal(err)
	}

	// new data for records
	r := data.New()

	// starting the application core
	core := app.New(c, r, s)

	// wrapper over dns lib
	d := dns.Server(c.UDPPort)
	go func() {
		if err := d.Run(c.IPV4); err != nil {
			log.Fatal("run dns service: ", err)
		}
	}()

	defer func() {
		if err := d.ShutDown(c.IPV4); err != nil {
			log.Println(err)
		}
	}()

	// for unmarshall from redis
	var m models.DNSEntry

	if len(res) > 0 {
		if err = json.Unmarshal(res, &m); err != nil {
			log.Fatal(err)
		}
		r.Set(m.Domain, &m)
	} else {
		if r, err = r.FetchCert(c); err != nil {
			log.Fatal(err)
		}
	}

	var swaggerSpec *loads.Document
	if swaggerSpec, err = loads.Analyzed(restapi.SwaggerJSON, ""); err != nil {
		log.Fatal("loads swagger spec", err)
	}

	api := operations.NewMdnsAPI(swaggerSpec)

	api.AddAddDNSEntryHandler = apiAdd.AddDNSEntryHandlerFunc(core.AddDNSEntryHandler)
	api.DeleteDeleteDNSEntryHandler = apiDelete.DeleteDNSEntryHandlerFunc(core.DeleteDNSEntryHandler)
	api.ShowListOneDNSEntryHandler = apiShow.ListOneDNSEntryHandlerFunc(core.ListOneDNSEntryHandler)
	api.ListShowDNSRecordsHandler = apiList.ShowDNSRecordsHandlerFunc(core.ShowDNSRecordsHandler)
	api.UpdateUpdateDNSEntryHandler = apiUpdate.UpdateDNSEntryHandlerFunc(core.UpdateDNSEntryHandler)

	server := restapi.NewServer(api)

	server.ConfigureAPI()

	var port int
	if port, err = strconv.Atoi(c.HTTPPort); err != nil {
		log.Fatal("can't convert port from string", err)
	}

	server.Port = port

	if err := server.Serve(); err != nil {
		log.Fatal("start server", err)
	}
}
