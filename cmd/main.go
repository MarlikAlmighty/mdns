package main

import (
	"github.com/MarlikAlmighty/mdns/internal/app"
	"github.com/MarlikAlmighty/mdns/internal/config"
	"github.com/MarlikAlmighty/mdns/internal/data"
	"github.com/MarlikAlmighty/mdns/internal/dns"
	"github.com/MarlikAlmighty/mdns/internal/gen/restapi"
	"github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations"
	apiAdd "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/add"
	apiDelete "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/delete"
	apiList "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/list"
	apiShow "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/show"
	apiUpdate "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/update"
	"github.com/go-openapi/loads"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {

	// get the configuration for the application through ENV
	cnf := config.New()
	if err := cnf.GetEnv(); err != nil {
		log.Fatalf("get env keys: %v\n", err)
	}

	// start new map for domains record
	dataMap := data.New()

	// starting the application core
	core := app.New(dataMap, cnf)

	// new dns server
	dnsServer := dns.New(dataMap, cnf)

	// start dns server
	dnsServer.Run()

	// for stopping
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// stop dns server
	go func() {
		for {
			select {
			case <-shutdown:
				if err := dnsServer.Close(); err != nil {
					log.Printf("%v\n", err)
				}
				return
			}
		}
	}()

	var (
		swaggerSpec *loads.Document
		err         error
	)

	if swaggerSpec, err = loads.Analyzed(restapi.SwaggerJSON, ""); err != nil {
		log.Fatalf("loads swagger spec %v\n", err)
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
	if port, err = strconv.Atoi(cnf.HTTPPort); err != nil {
		log.Fatalf("%v\n", err)
	}

	server.GracefulTimeout = 3 * time.Second
	server.Port = port
	server.Host = "127.0.0.1"

	if err = server.Serve(); err != nil {
		log.Fatalf("start rest api server %v\n", err)
	}
}
