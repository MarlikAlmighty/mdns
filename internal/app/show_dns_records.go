package app

import (
	apiList "github.com/MarlikAlmighty/mdns/internal/gen/restapi/operations/list"
	"github.com/go-openapi/runtime/middleware"
)

func (core *Core) ShowDNSRecordsHandler(params apiList.ShowDNSRecordsParams) middleware.Responder {
	return middleware.NotImplemented("operation list ShowDNSRecords has not yet been implemented")
}
