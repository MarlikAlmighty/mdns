package app

import (
	"github.com/MarlikAlmighty/mdns/internal/gen/models"
	apiAdd "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/add"
	"github.com/go-openapi/runtime/middleware"
)

func (core *Core) AddDNSEntryHandler(params apiAdd.AddDNSEntryParams) middleware.Responder {

	/*
		// TODO
		md, err := core.Resolver.FetchCert(params.Add.Domain, params.Add.IPV4)
		if err != nil {
			return apiAdd.NewAddDNSEntryBadRequest().WithPayload(&models.Answer{
				Code:    400,
				Message: err.Error(),
			})
		}
		log.Printf("CHALLENGE: %v\n", md.Acme)
	*/

	var (
		ipv6 string
		err  error
	)

	if ipv6, err = core.IPV4ToIPV6(params.Add.IPV4); err != nil {
		return apiAdd.NewAddDNSEntryBadRequest().WithPayload(&models.Answer{
			Code:    400,
			Message: "can't convert ipv4 to ipv6",
		})
	}

	params.Add.IPV6 = ipv6
	core.Resolver.Set(params.Add.Domain, params.Add)
	return apiAdd.NewAddDNSEntryOK().WithPayload(params.Add)
}
