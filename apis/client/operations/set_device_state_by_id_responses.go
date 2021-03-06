// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/leakingtapan/sonoff/apis/models"
)

// SetDeviceStateByIDReader is a Reader for the SetDeviceStateByID structure.
type SetDeviceStateByIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *SetDeviceStateByIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewSetDeviceStateByIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewSetDeviceStateByIDNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewSetDeviceStateByIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewSetDeviceStateByIDOK creates a SetDeviceStateByIDOK with default headers values
func NewSetDeviceStateByIDOK() *SetDeviceStateByIDOK {
	return &SetDeviceStateByIDOK{}
}

/*SetDeviceStateByIDOK handles this case with default header values.

successful operation
*/
type SetDeviceStateByIDOK struct {
	Payload *models.Device
}

func (o *SetDeviceStateByIDOK) Error() string {
	return fmt.Sprintf("[POST /devices/{deviceId}/{state}][%d] setDeviceStateByIdOK  %+v", 200, o.Payload)
}

func (o *SetDeviceStateByIDOK) GetPayload() *models.Device {
	return o.Payload
}

func (o *SetDeviceStateByIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Device)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewSetDeviceStateByIDNotFound creates a SetDeviceStateByIDNotFound with default headers values
func NewSetDeviceStateByIDNotFound() *SetDeviceStateByIDNotFound {
	return &SetDeviceStateByIDNotFound{}
}

/*SetDeviceStateByIDNotFound handles this case with default header values.

Device not found
*/
type SetDeviceStateByIDNotFound struct {
}

func (o *SetDeviceStateByIDNotFound) Error() string {
	return fmt.Sprintf("[POST /devices/{deviceId}/{state}][%d] setDeviceStateByIdNotFound ", 404)
}

func (o *SetDeviceStateByIDNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewSetDeviceStateByIDDefault creates a SetDeviceStateByIDDefault with default headers values
func NewSetDeviceStateByIDDefault(code int) *SetDeviceStateByIDDefault {
	return &SetDeviceStateByIDDefault{
		_statusCode: code,
	}
}

/*SetDeviceStateByIDDefault handles this case with default header values.

Internal Server Error
*/
type SetDeviceStateByIDDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the set device state by Id default response
func (o *SetDeviceStateByIDDefault) Code() int {
	return o._statusCode
}

func (o *SetDeviceStateByIDDefault) Error() string {
	return fmt.Sprintf("[POST /devices/{deviceId}/{state}][%d] setDeviceStateById default  %+v", o._statusCode, o.Payload)
}

func (o *SetDeviceStateByIDDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *SetDeviceStateByIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
