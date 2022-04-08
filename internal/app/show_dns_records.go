package app

import (
	apiList "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/list"
	"github.com/go-openapi/runtime/middleware"
)

func (core *Core) ShowDNSRecordsHandler(_ apiList.ShowDNSRecordsParams) middleware.Responder {
	return apiList.NewShowDNSRecordsOK().WithPayload(core.Resolver.GetMap())
}
