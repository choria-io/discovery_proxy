package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// DiscoveryRequest discovery request
// swagger:model discoveryRequest
type DiscoveryRequest struct {

	// agents
	Agents AgentsFilter `json:"agents"`

	// classes
	Classes ClassesFilter `json:"classes"`

	// collective
	Collective CollectiveFilter `json:"collective,omitempty"`

	// facts
	Facts FactsFilter `json:"facts"`

	// identities
	Identities IdentitiesFilter `json:"identities"`

	// node set
	NodeSet Word `json:"node_set,omitempty"`

	// PQL Query
	// Min Length: 1
	Query string `json:"query,omitempty"`
}

// Validate validates this discovery request
func (m *DiscoveryRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCollective(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateNodeSet(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateQuery(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DiscoveryRequest) validateCollective(formats strfmt.Registry) error {

	if swag.IsZero(m.Collective) { // not required
		return nil
	}

	if err := m.Collective.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("collective")
		}
		return err
	}

	return nil
}

func (m *DiscoveryRequest) validateNodeSet(formats strfmt.Registry) error {

	if swag.IsZero(m.NodeSet) { // not required
		return nil
	}

	if err := m.NodeSet.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("node_set")
		}
		return err
	}

	return nil
}

func (m *DiscoveryRequest) validateQuery(formats strfmt.Registry) error {

	if swag.IsZero(m.Query) { // not required
		return nil
	}

	if err := validate.MinLength("query", "body", string(m.Query), 1); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *DiscoveryRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DiscoveryRequest) UnmarshalBinary(b []byte) error {
	var res DiscoveryRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}