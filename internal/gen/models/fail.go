// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Fail fail
//
// swagger:model fail
type Fail struct {

	// code
	Code uint32 `json:"Code,omitempty"`

	// message
	Message string `json:"Message,omitempty"`
}

// Validate validates this fail
func (m *Fail) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this fail based on context it is used
func (m *Fail) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Fail) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Fail) UnmarshalBinary(b []byte) error {
	var res Fail
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
