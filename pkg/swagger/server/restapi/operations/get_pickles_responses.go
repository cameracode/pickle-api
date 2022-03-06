// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"pickle-api/pkg/swagger/server/models"
)

// GetPicklesOKCode is the HTTP code returned for type GetPicklesOK
const GetPicklesOKCode int = 200

/*GetPicklesOK Return the Pickles list.

swagger:response getPicklesOK
*/
type GetPicklesOK struct {

	/*
	  In: Body
	*/
	Payload []*models.Pickle `json:"body,omitempty"`
}

// NewGetPicklesOK creates GetPicklesOK with default headers values
func NewGetPicklesOK() *GetPicklesOK {

	return &GetPicklesOK{}
}

// WithPayload adds the payload to the get pickles o k response
func (o *GetPicklesOK) WithPayload(payload []*models.Pickle) *GetPicklesOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get pickles o k response
func (o *GetPicklesOK) SetPayload(payload []*models.Pickle) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPicklesOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]*models.Pickle, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}
