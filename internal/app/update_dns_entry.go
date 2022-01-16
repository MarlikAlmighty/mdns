package app

import (
	apiUpdate "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/update"
	"github.com/go-openapi/runtime/middleware"
)

func (core *Core) UpdateDNSEntryHandler(params apiUpdate.UpdateDNSEntryParams) middleware.Responder {
	return middleware.NotImplemented("operation update UpdateDNSEntry has not yet been implemented")
}
