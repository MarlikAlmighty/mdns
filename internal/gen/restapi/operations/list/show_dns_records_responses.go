// Code generated by go-swagger; DO NOT EDIT.

package list

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/MarlikAlmighty/mdns/internal/gen/models"
)

// ShowDNSRecordsOKCode is the HTTP code returned for type ShowDNSRecordsOK
const ShowDNSRecordsOKCode int = 200

/*
ShowDNSRecordsOK OK

swagger:response showDnsRecordsOK
*/
type ShowDNSRecordsOK struct {

	/*
	  In: Body
	*/
	Payload models.DNSRecords `json:"body,omitempty"`
}

// NewShowDNSRecordsOK creates ShowDNSRecordsOK with default headers values
func NewShowDNSRecordsOK() *ShowDNSRecordsOK {

	return &ShowDNSRecordsOK{}
}

// WithPayload adds the payload to the show Dns records o k response
func (o *ShowDNSRecordsOK) WithPayload(payload models.DNSRecords) *ShowDNSRecordsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the show Dns records o k response
func (o *ShowDNSRecordsOK) SetPayload(payload models.DNSRecords) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ShowDNSRecordsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty map
		payload = models.DNSRecords{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// ShowDNSRecordsBadRequestCode is the HTTP code returned for type ShowDNSRecordsBadRequest
const ShowDNSRecordsBadRequestCode int = 400

/*
ShowDNSRecordsBadRequest Bad request

swagger:response showDnsRecordsBadRequest
*/
type ShowDNSRecordsBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Answer `json:"body,omitempty"`
}

// NewShowDNSRecordsBadRequest creates ShowDNSRecordsBadRequest with default headers values
func NewShowDNSRecordsBadRequest() *ShowDNSRecordsBadRequest {

	return &ShowDNSRecordsBadRequest{}
}

// WithPayload adds the payload to the show Dns records bad request response
func (o *ShowDNSRecordsBadRequest) WithPayload(payload *models.Answer) *ShowDNSRecordsBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the show Dns records bad request response
func (o *ShowDNSRecordsBadRequest) SetPayload(payload *models.Answer) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ShowDNSRecordsBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
