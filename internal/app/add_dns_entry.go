package app

import (
	"github.com/MarlikAlmighty/mdns/internal/gen/models"
	apiAdd "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/add"
	"github.com/go-openapi/runtime/middleware"
	"log"
)

func (core *Core) AddDNSEntryHandler(params apiAdd.AddDNSEntryParams) middleware.Responder {

	// TODO enable delayed domain registration because it will take some time
	mp, err := core.Resolver.FetchCert(params.Add.Domain, params.Add.IPV4)
	if err != nil {
		return apiAdd.NewAddDNSEntryBadRequest().WithPayload(&models.Answer{
			Code:    400,
			Message: err.Error(),
		})
	}

	log.Printf("CHALLENGE: %v\n", mp.Acme)

	var ipv6 string

	if ipv6, err = core.IPV4ToIPV6(params.Add.IPV4); err != nil {
		return apiAdd.NewAddDNSEntryBadRequest().WithPayload(&models.Answer{
			Code:    400,
			Message: "can't convert ipv4 to ipv6",
		})
	}

	mp.IPV6 = ipv6

	if b := core.Resolver.Set(mp.Domain, mp); b == true {
		return apiAdd.NewAddDNSEntryBadRequest().WithPayload(&models.Answer{
			Code:    400,
			Message: "such the record " + params.Add.Domain + " exist",
		})
	}

	return apiAdd.NewAddDNSEntryOK().WithPayload(mp)
}
