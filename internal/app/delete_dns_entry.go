package app

import (
	"github.com/MarlikAlmighty/mdns/internal/gen/models"
	apiDelete "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/delete"
	"github.com/go-openapi/runtime/middleware"
)

func (core *Core) DeleteDNSEntryHandler(params apiDelete.DeleteDNSEntryParams) middleware.Responder {

	m := core.Resolver.Get(params.Delete.Domain)
	if m.Domain == "" {
		return apiDelete.NewDeleteDNSEntryBadRequest().WithPayload(&models.Answer{
			Code:    400,
			Message: "domain does not exist",
		})
	}

	core.Resolver.Delete(params.Delete.Domain)
	return apiDelete.NewDeleteDNSEntryOK().WithPayload(&models.Answer{
		Code:    200,
		Message: "OK",
	})
}
