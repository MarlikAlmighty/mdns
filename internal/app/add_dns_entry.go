package app

import (
	"crypto/rsa"
	"github.com/MarlikAlmighty/mdns/internal/config"
	"github.com/MarlikAlmighty/mdns/internal/gen/models"
	apiAdd "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/add"
	"github.com/go-openapi/runtime/middleware"
)

func (core *Core) AddDNSEntryHandler(params apiAdd.AddDNSEntryParams) middleware.Responder {

	var (
		ipv6 string
		err  error
	)

	if core.Config.(*config.Configuration).IPV6 {
		if ipv6, err = core.IPV4ToIPV6(params.Add.IPV4); err != nil {
			return apiAdd.NewAddDNSEntryBadRequest().WithPayload(&models.Answer{
				Code:    400,
				Message: "can't convert ipv4 to ipv6",
			})
		}
		params.Add.IPV6 = ipv6
	} else {
		params.Add.IPV6 = ""
	}

	var (
		privRSA *rsa.PrivateKey
		pubRSA  *rsa.PublicKey
	)

	if privRSA, pubRSA, err = core.Resolver.GenerateRsaKeyPair(); err != nil {
		return apiAdd.NewAddDNSEntryBadRequest().WithPayload(&models.Answer{
			Code:    400,
			Message: "can't generate rsa pair",
		})
	}

	var pubStr string
	if pubStr, err = core.Resolver.ExportRsaPublicKeyAsStr(pubRSA); err != nil {
		return apiAdd.NewAddDNSEntryBadRequest().WithPayload(&models.Answer{
			Code:    400,
			Message: "can't convert public cert to string",
		})
	}

	params.Add.DkimPrivateKey = core.Resolver.ExportRsaPrivateKeyAsStr(privRSA)
	params.Add.DkimPublicKey = pubStr
	params.Add.Acme = []string{""}
	core.Resolver.Set(params.Add.Domain, params.Add)
	return apiAdd.NewAddDNSEntryOK().WithPayload(params.Add)
}
