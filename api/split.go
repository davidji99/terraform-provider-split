package api

import (
	"fmt"
	"github.com/davidji99/simpleresty"
	"net/url"
)

// SplitsService handles communication with the segment related
// methods of the Split.io APIv2.
//
// Reference: https://docs.split.io/reference/splits-overview
type SplitsService service

// Split is a feature flag, toggle, or experiment.
type Split struct {
	ID                     *string             `json:"id"`
	Name                   *string             `json:"name"`
	Description            *string             `json:"description"`
	Tags                   []string            `json:"tags,omitempty"`
	CreationTime           *int64              `json:"creationTime"`
	RolloutStatusTimestamp *int64              `json:"rolloutStatusTimestamp"`
	TrafficType            *TrafficType        `json:"trafficType"`
	RolloutStatus          *SplitRolloutStatus `json:"rolloutStatus"`
}

// SplitRolloutStatus represents the rollout status.
type SplitRolloutStatus struct {
	ID   *string `json:"id"`
	Name *string `json:"name"`
}

// Splits represents all splits.
type Splits struct {
	Objects []*Split `json:"objects"`
	GenericListResult
}

// SplitCreateRequest represents a request to create a split.
type SplitCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// SplitUpdateRequest represents a request to update a split.
type SplitUpdateRequest struct {
	Description string `json:"description"`
}

// List all splits.
//
// Reference: https://docs.split.io/reference/list-splits
func (s *SplitsService) List(workspaceId string, opts ...interface{}) (*Splits, *simpleresty.Response, error) {
	var result Splits
	urlStr, err := s.client.http.RequestURLWithQueryParams(fmt.Sprintf("/splits/ws/%s", workspaceId), opts...)
	if err != nil {
		return nil, nil, err
	}

	// Execute the request
	response, getErr := s.client.http.Get(urlStr, &result, nil)

	return &result, response, getErr
}

// Get a single split.
//
// splitId can be either the name or the UUID.
//
// Reference: https://docs.split.io/reference/get-split
func (s *SplitsService) Get(workspaceId, splitId string) (*Split, *simpleresty.Response, error) {
	var result Split
	urlStr := s.client.http.RequestURL("/splits/ws/%s/%s", workspaceId, splitId)

	// Execute the request
	response, getErr := s.client.http.Get(urlStr, &result, nil)

	return &result, response, getErr
}

// Create a single split.
//
// Reference: https://docs.split.io/reference/create-split
func (s *SplitsService) Create(workspaceId, trafficTypeId string, opts *SplitCreateRequest) (*Split, *simpleresty.Response, error) {
	var result Split
	urlStr := s.client.http.RequestURL("/splits/ws/%s/trafficTypes/%s", workspaceId, trafficTypeId)

	// Execute the request
	response, createErr := s.client.http.Post(urlStr, &result, opts)

	return &result, response, createErr
}

// UpdateDescription of an existing split.
//
// Reference: https://docs.split.io/reference/update-split-description
func (s *SplitsService) UpdateDescription(workspaceId, splitName, description string) (*Split, *simpleresty.Response, error) {
	var result Split

	splitNameEncoded := url.QueryEscape(splitName)
	urlStr := s.client.http.RequestURL("/splits/ws/%s/%s/updateDescription", workspaceId, splitNameEncoded)

	// Execute the request
	response, updateErr := s.client.http.Put(urlStr, &result, description)

	return &result, response, updateErr
}

// Delete a single split.
//
// This will automatically unconfigure the Split Definition from all environments. Returns `true` in the response body.
//
// Split name is required, not the split UUID.
//
// Reference: https://docs.split.io/reference/delete-split
func (s *SplitsService) Delete(workspaceId, splitName string) (*simpleresty.Response, error) {
	splitNameEncoded := url.QueryEscape(splitName)
	urlStr := s.client.http.RequestURL("/splits/ws/%s/%s", workspaceId, splitNameEncoded)

	// Execute the request
	response, createErr := s.client.http.Delete(urlStr, nil, nil)

	return response, createErr
}
