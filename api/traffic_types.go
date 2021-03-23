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
	Name               *string    `json:"name"`
	Type               *string    `json:"type"`
	ID                 *string    `json:"id"`
	DisplayAttributeID *string    `json:"displayAttributeId"`
	Workspace          *Workspace `json:"workspace"`
}

// List all traffic types.
//
// Reference: https://docs.split.io/reference#get-traffic-types
func (t *TrafficTypesService) List(workspaceID string) ([]*TrafficType, *simpleresty.Response, error) {
	var result []*TrafficType
	urlStr := t.client.http.RequestURL("/trafficTypes/ws/%s", workspaceID)

	// Execute the request
	response, getErr := t.client.http.Get(urlStr, &result, nil)

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
