// Code generated by go-swagger; DO NOT EDIT.

package delete

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// DeleteDNSEntryHandlerFunc turns a function with the right signature into a delete dns entry handler
type DeleteDNSEntryHandlerFunc func(DeleteDNSEntryParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteDNSEntryHandlerFunc) Handle(params DeleteDNSEntryParams) middleware.Responder {
	return fn(params)
}

// DeleteDNSEntryHandler interface for that can handle valid delete dns entry params
type DeleteDNSEntryHandler interface {
	Handle(DeleteDNSEntryParams) middleware.Responder
}

// NewDeleteDNSEntry creates a new http.Handler for the delete dns entry operation
func NewDeleteDNSEntry(ctx *middleware.Context, handler DeleteDNSEntryHandler) *DeleteDNSEntry {
	return &DeleteDNSEntry{Context: ctx, Handler: handler}
}

/* DeleteDNSEntry swagger:route DELETE /dns delete deleteDnsEntry

Delete dns entry

*/
type DeleteDNSEntry struct {
	Context *middleware.Context
	Handler DeleteDNSEntryHandler
}

func (o *DeleteDNSEntry) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewDeleteDNSEntryParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
