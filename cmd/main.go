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

	apiCerts "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/certs"

	apiShow "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/show"

	apiList "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/list"

	apiUpdate "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/update"

	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-openapi/loads"
)

func main() {

	// get the configuration for the application through ENV
	cnf := config.New()
	if err := cnf.GetEnv(); err != nil {
		log.Fatalf("get environment keys: %v\n", err)
	}

	// start new map for domains record
	dataMap := data.New()

	// new dns server
	dnsServer := dns.New(cnf.NameServers, cnf.DnsHost, dataMap)

	// start dns server
	if err := dnsServer.Run(); err != nil {
		log.Fatalf("run dns service: %v", err)
	}

	// for stopping
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// stop dns server
	go func() {
		for {
			select {
			case <-shutdown:
				if err := dnsServer.Close(); err != nil {
					log.Printf("stop dns service: %v", err)
					return
				}
			}
		}
	}()

	// starting the application core
	core := app.New(dataMap, cnf)

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
	api.CertsFetchCertsHandler = apiCerts.FetchCertsHandlerFunc(core.FetchCertsHandler)
	api.ShowListOneDNSEntryHandler = apiShow.ListOneDNSEntryHandlerFunc(core.ListOneDNSEntryHandler)
	api.ListShowDNSRecordsHandler = apiList.ShowDNSRecordsHandlerFunc(core.ShowDNSRecordsHandler)
	api.UpdateUpdateDNSEntryHandler = apiUpdate.UpdateDNSEntryHandlerFunc(core.UpdateDNSEntryHandler)

	server := restapi.NewServer(api)

	server.ConfigureAPI()

	var port int
	if port, err = strconv.Atoi(cnf.HTTPPort); err != nil {
		log.Fatal("can't convert port from string", err)
	}

	server.GracefulTimeout = 3 * time.Second
	server.Port = port
	server.Host = "127.0.0.1"

	if err = server.Serve(); err != nil {
		log.Fatal("start server", err)
	}
}
