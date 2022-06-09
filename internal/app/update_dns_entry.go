package app

import (
	"github.com/MarlikAlmighty/mdns/internal/gen/models"
	apiUpdate "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/update"
	"github.com/go-openapi/runtime/middleware"
)

func (core *Core) UpdateDNSEntryHandler(params apiUpdate.UpdateDNSEntryParams) middleware.Responder {

	m := core.Resolver.Get(params.Update.Domain)
	if m.Domain == "" {
		return apiUpdate.NewUpdateDNSEntryBadRequest().WithPayload(&models.Answer{
			Code:    400,
			Message: "domain does not exist",
		})
	}

	m.Domain = params.Update.Domain
	m.IPV4 = params.Update.IPV4
	m.Ips = params.Update.Ips
	core.Resolver.Set(m.Domain, m)
	return apiUpdate.NewUpdateDNSEntryOK().WithPayload(m)
}
