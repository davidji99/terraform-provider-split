package api

import (
	"fmt"

	"github.com/davidji99/simpleresty"
)

// TrafficTypesService handles communication with the traffic types related
// methods of the Split.io APIv2.
//
// Reference: https://docs.split.io/reference#traffic-types-overview
type TrafficTypesService service

type TrafficType struct {
	Name               *string    `json:"name,omitempty"`
	Type               *string    `json:"type,omitempty"`
	ID                 *string    `json:"id,omitempty"`
	DisplayAttributeID *string    `json:"displayAttributeId,omitempty"`
	Workspace          *Workspace `json:"workspace,omitempty"`
}

type TrafficTypeRequest struct {
	// Name must begin with a letter (a-z or A-Z) and can only contain letters, numbers, hyphens and underscores
	// (a-z, A-Z, 0-9, -, _).
	Name string `json:"name"`
}

// List all traffic types.
//
// Reference: https://docs.split.io/reference#get-traffic-types
func (t *TrafficTypesService) List(workspaceID string) ([]*TrafficType, *simpleresty.Response, error) {
	var result []*TrafficType
	urlStr := t.client.http.RequestURL("/trafficTypes/ws/%s", workspaceID)
	response, getErr := t.client.get(urlStr, &result, nil)

	return result, response, getErr
}

// FindByID retrieves a traffic type by its ID.
func (t *TrafficTypesService) FindByID(workspaceID, trafficTypeID string) (*TrafficType, *simpleresty.Response, error) {
	trafficTypes, response, listErr := t.List(workspaceID)
	if listErr != nil {
		return nil, response, listErr
	}

	for _, t := range trafficTypes {
		if t.GetID() == trafficTypeID {
			return t, nil, nil
		}
	}

	return nil, nil, fmt.Errorf("traffic type [%s] not found", trafficTypeID)
}

// FindByName retrieves a traffic type by its name.
func (t *TrafficTypesService) FindByName(workspaceID, trafficTypeName string) (*TrafficType, *simpleresty.Response, error) {
	trafficTypes, response, listErr := t.List(workspaceID)
	if listErr != nil {
		return nil, response, listErr
	}

	for _, t := range trafficTypes {
		if t.GetName() == trafficTypeName {
			return t, nil, nil
		}
	}

	return nil, nil, fmt.Errorf("traffic type [%s] not found", trafficTypeName)
}

// Create a traffic type.
//
// Reference: https://docs.split.io/reference/create-traffic-types
func (t *TrafficTypesService) Create(workspaceID string, opts *TrafficTypeRequest) (*TrafficType, *simpleresty.Response, error) {
	var result TrafficType
	urlStr := t.client.http.RequestURL("/trafficTypes/ws/%s", workspaceID)

	// Execute the request
	response, createErr := t.client.post(urlStr, &result, opts)

	return &result, response, createErr
}

// Delete a traffic type.
//
// Reference: https://docs.split.io/reference/delete-trafic-type
func (t *TrafficTypesService) Delete(trafficTypeID string) (*simpleresty.Response, error) {
	urlStr := t.client.http.RequestURL("/trafficTypes/%s", trafficTypeID)

	// Execute the request
	response, deleteErr := t.client.delete(urlStr, nil, nil)

	return response, deleteErr
}
