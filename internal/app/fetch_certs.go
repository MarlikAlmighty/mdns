package app

import (
	"github.com/MarlikAlmighty/mdns/internal/gen/models"
	apiAdd "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/add"
	apiCerts "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/certs"
	"github.com/go-openapi/runtime/middleware"
)

func (core *Core) FetchCertsHandler(params apiCerts.FetchCertsParams) middleware.Responder {

	var (
		md  *models.DNSEntry
		err error
	)

	if md, err = core.Resolver.FetchCert(params.Certs.Domain, params.Certs.IPV4); err != nil {
		return apiAdd.NewAddDNSEntryBadRequest().WithPayload(&models.Answer{
			Code:    400,
			Message: err.Error(),
		})
	}

	return apiCerts.NewFetchCertsOK().WithPayload(md)
}
