package app

import (
	"github.com/MarlikAlmighty/mdns/internal/gen/models"
	apiAdd "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/add"
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
	m.Ipv4s = params.Update.Ipv4s
	m.Acme = params.Update.Acme

	var (
		ipv6 string
		err  error
	)

	m.Ipv6s = []string{}

	for _, v := range params.Update.Ipv4s {
		ipv6, err = core.IPV4ToIPV6(v)
		if err != nil {
			return apiAdd.NewAddDNSEntryBadRequest().WithPayload(&models.Answer{
				Code:    400,
				Message: "can't convert ipv4 to ipv6",
			})
		}
		m.Ipv6s = append(m.Ipv6s, ipv6)
	}

	core.Resolver.Set(m.Domain, m)
	return apiUpdate.NewUpdateDNSEntryOK().WithPayload(m)
}
