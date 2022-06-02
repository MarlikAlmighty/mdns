// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations"
	"github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/add"
	"github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/certs"
	"github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/delete"
	"github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/list"
	"github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/show"
	"github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/update"
)

//go:generate swagger generate server --target ../../gen --name Mdns --spec ../../../swagger-api/swagger.yml --template-dir ./swagger-templates --principal interface{}

func configureFlags(api *operations.MdnsAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.MdnsAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.AddAddDNSEntryHandler == nil {
		api.AddAddDNSEntryHandler = add.AddDNSEntryHandlerFunc(func(params add.AddDNSEntryParams) middleware.Responder {
			return middleware.NotImplemented("operation add.AddDNSEntry has not yet been implemented")
		})
	}
	if api.DeleteDeleteDNSEntryHandler == nil {
		api.DeleteDeleteDNSEntryHandler = delete.DeleteDNSEntryHandlerFunc(func(params delete.DeleteDNSEntryParams) middleware.Responder {
			return middleware.NotImplemented("operation delete.DeleteDNSEntry has not yet been implemented")
		})
	}
	if api.CertsFetchCertsHandler == nil {
		api.CertsFetchCertsHandler = certs.FetchCertsHandlerFunc(func(params certs.FetchCertsParams) middleware.Responder {
			return middleware.NotImplemented("operation certs.FetchCerts has not yet been implemented")
		})
	}
	if api.ShowListOneDNSEntryHandler == nil {
		api.ShowListOneDNSEntryHandler = show.ListOneDNSEntryHandlerFunc(func(params show.ListOneDNSEntryParams) middleware.Responder {
			return middleware.NotImplemented("operation show.ListOneDNSEntry has not yet been implemented")
		})
	}
	if api.ListShowDNSRecordsHandler == nil {
		api.ListShowDNSRecordsHandler = list.ShowDNSRecordsHandlerFunc(func(params list.ShowDNSRecordsParams) middleware.Responder {
			return middleware.NotImplemented("operation list.ShowDNSRecords has not yet been implemented")
		})
	}
	if api.UpdateUpdateDNSEntryHandler == nil {
		api.UpdateUpdateDNSEntryHandler = update.UpdateDNSEntryHandlerFunc(func(params update.UpdateDNSEntryParams) middleware.Responder {
			return middleware.NotImplemented("operation update.UpdateDNSEntry has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
