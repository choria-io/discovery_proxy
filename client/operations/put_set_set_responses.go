package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/choria-io/discovery_proxy/models"
)

// PutSetSetReader is a Reader for the PutSetSet structure.
type PutSetSetReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PutSetSetReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewPutSetSetOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 404:
		result := NewPutSetSetNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 500:
		result := NewPutSetSetInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewPutSetSetOK creates a PutSetSetOK with default headers values
func NewPutSetSetOK() *PutSetSetOK {
	return &PutSetSetOK{}
}

/*PutSetSetOK handles this case with default header values.

Basic successful request
*/
type PutSetSetOK struct {
	Payload *models.SuccessModel
}

func (o *PutSetSetOK) Error() string {
	return fmt.Sprintf("[PUT /set/{set}][%d] putSetSetOK  %+v", 200, o.Payload)
}

func (o *PutSetSetOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.SuccessModel)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPutSetSetNotFound creates a PutSetSetNotFound with default headers values
func NewPutSetSetNotFound() *PutSetSetNotFound {
	return &PutSetSetNotFound{}
}

/*PutSetSetNotFound handles this case with default header values.

Not found
*/
type PutSetSetNotFound struct {
}

func (o *PutSetSetNotFound) Error() string {
	return fmt.Sprintf("[PUT /set/{set}][%d] putSetSetNotFound ", 404)
}

func (o *PutSetSetNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewPutSetSetInternalServerError creates a PutSetSetInternalServerError with default headers values
func NewPutSetSetInternalServerError() *PutSetSetInternalServerError {
	return &PutSetSetInternalServerError{}
}

/*PutSetSetInternalServerError handles this case with default header values.

Standard Error Format
*/
type PutSetSetInternalServerError struct {
	Payload *models.ErrorModel
}

func (o *PutSetSetInternalServerError) Error() string {
	return fmt.Sprintf("[PUT /set/{set}][%d] putSetSetInternalServerError  %+v", 500, o.Payload)
}

func (o *PutSetSetInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorModel)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
