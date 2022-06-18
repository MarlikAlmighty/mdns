package app

import (
	"crypto/rsa"
	"github.com/MarlikAlmighty/mdns/internal/gen/models"
	apiAdd "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/add"
	"github.com/go-openapi/runtime/middleware"
)

func (core *Core) AddDNSEntryHandler(params apiAdd.AddDNSEntryParams) middleware.Responder {

	md := core.Resolver.Get(params.Add.Domain)

	var (
		ipv6 string
		err  error
	)

	md.Ipv6s = []string{}
	for _, v := range params.Add.Ipv4s {
		ipv6, err = core.IPV4ToIPV6(v)
		if err != nil {
			return apiAdd.NewAddDNSEntryBadRequest().WithPayload(&models.Answer{
				Code:    400,
				Message: "can't convert ipv4 to ipv6",
			})
		}
		md.Ipv6s = append(md.Ipv6s, ipv6)
	}

	var (
		privRSA *rsa.PrivateKey
		pubRSA  *rsa.PublicKey
	)

	if privRSA, pubRSA, err = core.GenerateRsaKeyPair(); err != nil {
		return apiAdd.NewAddDNSEntryBadRequest().WithPayload(&models.Answer{
			Code:    400,
			Message: "can't generate rsa pair",
		})
	}

	var pubStr string
	if pubStr, err = core.ExportRsaPublicKeyAsStr(pubRSA); err != nil {
		return apiAdd.NewAddDNSEntryBadRequest().WithPayload(&models.Answer{
			Code:    400,
			Message: "can't convert public cert to string",
		})
	}

	if md.Domain == "" || len(md.Ipv4s) == 0 {
		md.Domain = params.Add.Domain
		md.Ipv4s = params.Add.Ipv4s
	}

	md.DkimPrivateKey = core.ExportRsaPrivateKeyAsStr(privRSA)
	md.DkimPublicKey = pubStr
	md.Acme = []string{""}
	core.Resolver.Set(md.Domain, md)
	return apiAdd.NewAddDNSEntryOK().WithPayload(md)
}
