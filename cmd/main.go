package main

import (
	"fmt"
	"github.com/MarlikAlmighty/mdns/internal/app"
	"github.com/MarlikAlmighty/mdns/internal/config"
	"github.com/MarlikAlmighty/mdns/internal/data"
	"github.com/MarlikAlmighty/mdns/internal/gen/restapi"
	"github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations"
	apiAdd "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/add"
	apiDelete "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/delete"
	apiList "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/list"
	apiShow "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/show"
	"github.com/MarlikAlmighty/mdns/internal/wrapper"

	apiUpdate "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/update"

	"log"
	"strconv"

	"github.com/go-openapi/loads"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered func, we got panic:", r)
		}
	}()

	c := config.New()

	if err := c.GetEnv(); err != nil {
		log.Fatal("get environment keys: ", err)
	}

	// TODO SQLite

	core := app.New(c)

	if err := core.GenerateCerts(); err != nil {
		log.Fatal("generating certs: ", err)
	}

	if c.IPV6 == "" {
		ipv6, err := core.IPV4ToIPV6(c.IPV4)
		if err != nil {
			log.Fatal("convert ipv4 to ipv6: ", err)
		}
		c.IPV6 = ipv6
	}

	if crt, key, err := core.CertsFromFile(); err != nil {
		log.Fatal("get certs from file: ", err)
	} else {
		c.PublicKey = crt
		c.PrivateKey = key
	}

	// TODO Change me
	_ = data.NewResolvedData(c)

	// wrapper over dns lib
	r := wrapper.New()
	if err := r.Run(); err != nil {
		log.Fatal("run dns service: ", err)
	}
	defer r.ShutDown()

	var (
		swaggerSpec *loads.Document
		err         error
	)

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
