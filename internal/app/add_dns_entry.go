package app

import (
	apiAdd "mdns/internal/gen/restapi/operations/add"

	"github.com/go-openapi/runtime/middleware"
)

func (core *Core) AddDNSEntryHandler(params apiAdd.AddDNSEntryParams) middleware.Responder {
	return middleware.NotImplemented("operation add AddDNSEntry has not yet been implemented")
}
