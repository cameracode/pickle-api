// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/runtime"
)

// GetPickleNameOKCode is the HTTP code returned for type GetPickleNameOK
const GetPickleNameOKCode int = 200

/*GetPickleNameOK Returns the Pickle Rick

swagger:response getPickleNameOK
*/
type GetPickleNameOK struct {

	/*
	  In: Body
	*/
	Payload io.ReadCloser `json:"body,omitempty"`
}

// NewGetPickleNameOK creates GetPickleNameOK with default headers values
func NewGetPickleNameOK() *GetPickleNameOK {

	return &GetPickleNameOK{}
}

// WithPayload adds the payload to the get pickle name o k response
func (o *GetPickleNameOK) WithPayload(payload io.ReadCloser) *GetPickleNameOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get pickle name o k response
func (o *GetPickleNameOK) SetPayload(payload io.ReadCloser) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPickleNameOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}