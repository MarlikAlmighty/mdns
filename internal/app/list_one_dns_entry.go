package app

import (
	apiShow "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/show"
	"github.com/go-openapi/runtime/middleware"
)

func (core *Core) ListOneDNSEntryHandler(params apiShow.ListOneDNSEntryParams) middleware.Responder {
	return apiShow.NewListOneDNSEntryOK().WithPayload(core.Resolver.Get(params.ID))
}
