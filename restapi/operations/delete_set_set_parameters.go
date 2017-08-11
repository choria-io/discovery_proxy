// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"
)

// NewDeleteSetSetParams creates a new DeleteSetSetParams object
// with the default values initialized.
func NewDeleteSetSetParams() DeleteSetSetParams {
	var ()
	return DeleteSetSetParams{}
}

// DeleteSetSetParams contains all the bound params for the delete set set operation
// typically these are obtained from a http.Request
//
// swagger:parameters DeleteSetSet
type DeleteSetSetParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request

	/*Node set to delete
	  Required: true
	  In: path
	*/
	Set string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *DeleteSetSetParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	o.HTTPRequest = r

	rSet, rhkSet, _ := route.Params.GetOK("set")
	if err := o.bindSet(rSet, rhkSet, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *DeleteSetSetParams) bindSet(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	o.Set = raw

	return nil
}
