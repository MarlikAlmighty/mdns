package main

import (
	"github.com/MarlikAlmighty/mdns/internal/app"
	"github.com/MarlikAlmighty/mdns/internal/config"
	"github.com/MarlikAlmighty/mdns/internal/data"
	"github.com/MarlikAlmighty/mdns/internal/dns"
	"github.com/MarlikAlmighty/mdns/internal/dump"
	"github.com/MarlikAlmighty/mdns/internal/gen/restapi"
	"github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations"
	apiAdd "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/add"

	apiDelete "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/delete"

	apiShow "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/show"

	apiList "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/list"

	apiUpdate "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/update"

	"log"
	"strconv"
	"time"

	"github.com/go-openapi/loads"
)

func main() {

	// for to catch panic
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered func, we got panic:", r)
		}
	}()

	// get the configuration for the application through ENV
	cnf := config.New()
	if err := cnf.GetEnv(); err != nil {
		log.Fatal("get environment keys: ", err)
	}

	// start new map for domains record
	dataMap := data.New()

	// TODO run a redis container yourself
	// start store
	redisClient, err := dump.New(cnf.RedisUrl, dataMap)
	if err != nil {
		log.Fatal("start redis: ", err)
	}

	// stop store
	defer func() {
		if err = redisClient.Shutdown(cnf.RedisKey); err != nil {
			log.Printf("stop redis client: %v", err)
		}
	}()

	// if we have value in dump, adding it to map
	if err = redisClient.InitMaps(cnf.RedisKey); err != nil {
		log.Fatal("redis init map: ", err)
	}

	// start dns server
	dnsServer := dns.New(cnf.UDPPort, dataMap)
	go func() {
		if err := dnsServer.Run(); err != nil {
			log.Fatal("run dns service: ", err)
		}
	}()

	// stop dns server
	defer func() {
		if err = dnsServer.Close(); err != nil {
			log.Printf("stop dns service: %v", err)
		}
	}()

	// starting the application core
	core := app.New(dataMap)

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

	server.GracefulTimeout = 1 * time.Second

	var port int
	if port, err = strconv.Atoi(cnf.HTTPPort); err != nil {
		log.Fatal("can't convert port from string", err)
	}

	server.Port = port

	if err := server.Serve(); err != nil {
		log.Fatal("start server", err)
	}
}
