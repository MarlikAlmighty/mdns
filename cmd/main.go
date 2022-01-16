package main

import (
	"github.com/MarlikAlmighty/mdns/internal/app"
	"github.com/MarlikAlmighty/mdns/internal/config"
	"github.com/MarlikAlmighty/mdns/internal/gen/restapi"
	"github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations"
	apiAdd "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/add"

	apiDelete "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/delete"

	apiShow "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/show"

	apiList "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/list"

	apiUpdate "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/update"

	"github.com/go-openapi/loads"
	"go.uber.org/zap"
)

func main() {

	var (
		l   *zap.Logger
		err error
	)

	if l, err = zap.NewDevelopment(); err != nil {
		panic(err)
	}

	var c *config.Configuration
	if c, err = config.New(); err != nil {
		l.Fatal("get environment keys", zap.Error(err))
	}

	core := app.New(c, l)

	var swaggerSpec *loads.Document
	if swaggerSpec, err = loads.Analyzed(restapi.SwaggerJSON, ""); err != nil {
		l.Fatal("loads swagger spec", zap.Error(err))
	}

	api := operations.NewMdnsAPI(swaggerSpec)

	api.AddAddDNSEntryHandler = apiAdd.AddDNSEntryHandlerFunc(core.AddDNSEntryHandler)
	api.DeleteDeleteDNSEntryHandler = apiDelete.DeleteDNSEntryHandlerFunc(core.DeleteDNSEntryHandler)
	api.ShowListOneDNSEntryHandler = apiShow.ListOneDNSEntryHandlerFunc(core.ListOneDNSEntryHandler)
	api.ListShowDNSRecordsHandler = apiList.ShowDNSRecordsHandlerFunc(core.ShowDNSRecordsHandler)
	api.UpdateUpdateDNSEntryHandler = apiUpdate.UpdateDNSEntryHandlerFunc(core.UpdateDNSEntryHandler)

	server := restapi.NewServer(api)

	server.ConfigureAPI()

	server.Port = int(c.HTTPPort)

	if err := server.Serve(); err != nil {
		l.Fatal("start server", zap.Error(err))
	}

	

}
