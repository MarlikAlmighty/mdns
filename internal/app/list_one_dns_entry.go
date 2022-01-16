package app

import (
	apiShow "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/show"
	"github.com/go-openapi/runtime/middleware"
)

func (core *Core) ListOneDNSEntryHandler(params apiShow.ListOneDNSEntryParams) middleware.Responder {
	return middleware.NotImplemented("operation show ListOneDNSEntry has not yet been implemented")
}
