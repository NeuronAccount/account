// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// SmsSignupOKCode is the HTTP code returned for type SmsSignupOK
const SmsSignupOKCode int = 200

/*SmsSignupOK ok

swagger:response smsSignupOK
*/
type SmsSignupOK struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewSmsSignupOK creates SmsSignupOK with default headers values
func NewSmsSignupOK() *SmsSignupOK {

	return &SmsSignupOK{}
}

// WithPayload adds the payload to the sms signup o k response
func (o *SmsSignupOK) WithPayload(payload string) *SmsSignupOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the sms signup o k response
func (o *SmsSignupOK) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SmsSignupOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
