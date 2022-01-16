package app

import (
	apiDelete "mdns/internal/gen/restapi/operations/delete"

	"github.com/go-openapi/runtime/middleware"
)

func (core *Core) DeleteDNSEntryHandler(params apiDelete.DeleteDNSEntryParams) middleware.Responder {
	return middleware.NotImplemented("operation delete DeleteDNSEntry has not yet been implemented")
}
