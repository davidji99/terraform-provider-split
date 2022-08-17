package api

import (
	"fmt"

	"github.com/davidji99/simpleresty"
	"net/url"
)

// AttributesService handles communication with the attributes related
// methods of the Split.io APIv2.
//
// Reference: https://docs.split.io/reference/attributes-overview
type AttributesService service

// Attribute represents an attribute in Split.
type Attribute struct {
	ID              *string  `json:"id"` // this is different from the usually computed ID.
	OrganizationId  *string  `json:"organizationId"`
	TrafficTypeID   *string  `json:"trafficTypeId"`
	DisplayName     *string  `json:"displayName"`
	Description     *string  `json:"description"`
	DataType        *string  `json:"dataType"` // (Optional) The data type of the attribute used for display formatting, defaults to displaying the raw string. Must be one of: null, "string", "datetime", "number", "set"
	IsSearchable    *bool    `json:"isSearchable"`
	SuggestedValues []string `json:"suggestedValues"`
}

// AttributeRequest represents a request to create an attribute.
type AttributeRequest struct {
	Identifier      *string  `json:"id,omitempty"`
	DisplayName     *string  `json:"displayName,omitempty"`
	Description     *string  `json:"description,omitempty"`
	TrafficTypeID   *string  `json:"trafficTypeId,omitempty"`
	IsSearchable    *bool    `json:"isSearchable,omitempty"`
	DataType        *string  `json:"dataType,omitempty"`
	SuggestedValues []string `json:"suggestedValues,omitempty"`
}

// List all attributes for a traffic type.
//
// Reference: https://docs.split.io/reference/get-attributes
func (a *AttributesService) List(workspaceID, trafficTypeID string) ([]*Attribute, *simpleresty.Response, error) {
	var result []*Attribute
	urlStr := a.client.http.RequestURL("/schema/ws/%s/trafficTypes/%s", workspaceID, trafficTypeID)

	response, listErr := a.client.get(urlStr, &result, nil)

	return result, response, listErr
}

// FindByID retrieves an attribute by its ID.
//
// This is a helper method as it is not possible to retrieve a single attribute.
func (a *AttributesService) FindByID(workspaceID, trafficTypeID, attributeID string) (*Attribute, *simpleresty.Response, error) {
	attributes, listErr, listResponse := a.List(workspaceID, trafficTypeID)
	if listErr != nil {
		return nil, listErr, listResponse
	}

	for _, a := range attributes {
		if a.GetID() == attributeID {
			return a, nil, nil
		}
	}

	return nil, nil, fmt.Errorf("could not find attribute [%s]", attributeID)
}

// Create an attribute.
//
// Reference: https://docs.split.io/reference/save-attribute
func (a *AttributesService) Create(workspaceID, trafficTypeID string, opts *AttributeRequest) (*Attribute, *simpleresty.Response, error) {
	var result Attribute
	urlStr := a.client.http.RequestURL("/schema/ws/%s/trafficTypes/%s", workspaceID, trafficTypeID)

	// Execute the request
	response, createErr := a.client.post(urlStr, &result, opts)

	return &result, response, createErr
}

// Update an attribute.
//
// Reference: https://docs.split.io/reference/update-attribute
func (a *AttributesService) Update(workspaceID, trafficTypeID, attributeID string, opts *AttributeRequest) (*Attribute, *simpleresty.Response, error) {
	var result Attribute
	attributeIdEncoded := url.QueryEscape(attributeID)
	urlStr := a.client.http.RequestURL("/schema/ws/%s/trafficTypes/%s/%s", workspaceID, trafficTypeID, attributeIdEncoded)

	// Execute the request
	response, createErr := a.client.patch(urlStr, &result, opts)

	return &result, response, createErr
}

// Delete an attribute.
//
// Reference: https://docs.split.io/reference/delete-attribute
func (a *AttributesService) Delete(workspaceID, trafficTypeID, attributeID string) (*simpleresty.Response, error) {
	attributeIdEncoded := url.QueryEscape(attributeID)
	urlStr := a.client.http.RequestURL("/schema/ws/%s/trafficTypes/%s/%s", workspaceID, trafficTypeID, attributeIdEncoded)

	// Execute the request
	response, deleteErr := a.client.delete(urlStr, nil, nil)

	return response, deleteErr
}
