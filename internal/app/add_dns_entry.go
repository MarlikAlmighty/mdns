package app

import (
	"github.com/MarlikAlmighty/mdns/internal/gen/models"
	apiAdd "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/add"
	"github.com/go-openapi/runtime/middleware"
)

func (core *Core) AddDNSEntryHandler(params apiAdd.AddDNSEntryParams) middleware.Responder {
	ipv6, err := core.IPV4ToIPV6(params.Add.IPV4)
	if err != nil {
		return apiAdd.NewAddDNSEntryBadRequest().WithPayload(&models.Answer{
			Code:    400,
			Message: "can't convert ipv4 to ipv6",
		})
	}
	m := &models.DNSEntry{}
	m.Domain = params.Add.Domain
	m.IPV4 = params.Add.IPV4
	m.IPV6 = ipv6
	// TODO enable pending domain registration
	/*
		mp, err := core.Resolver.FetchCert(core.Config.(*config.Configuration))
		if err != nil {
			return apiAdd.NewAddDNSEntryBadRequest().WithPayload(&models.Fail{
				Code:    400,
				Message: err.Error(),
			})
		}
	*/
	core.Resolver.Set(m.Domain, m)
	return apiAdd.NewAddDNSEntryOK().WithPayload(m)
}
