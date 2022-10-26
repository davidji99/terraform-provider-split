package api

import (
	"github.com/davidji99/simpleresty"
)

// GroupsService handles communication with the group related
// methods of the Split.io APIv2.
//
// Reference: https://docs.split.io/reference#groups-overview
type GroupsService service

// Group represents a group in Split.
type Group struct {
	ID          *string `json:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Type        *string `json:"type"`
}

// GroupListResult
type GroupListResult struct {
	Data           []*Group `json:"objects"`
	NextMarker     *string  `json:"nextMarker,omitempty"`
	PreviousMarker *string  `json:"previousMarker,omitempty"`
	Limit          *int     `json:"limit"`
	Count          *int     `json:"count"`
}

// GroupListOpts
type GroupListOpts struct {
	// 1-200 are the potential values. Default=50
	Limit int `url:"limit,omitempty"`
}

// GroupRequest
type GroupRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// List all active Groups in the organization
//
// Reference: https://docs.split.io/reference#list-groups
func (g *GroupsService) List(opts *GroupListOpts) (*GroupListResult, *simpleresty.Response, error) {
	var result GroupListResult
	//Here I have a problem with the RequestURLWithParams, TBD
	urlStr, urlStrErr := g.client.http.RequestURLWithQueryParams("/groups", opts)
	if urlStrErr != nil {
		return nil, nil, urlStrErr
	}

	// Execute the request
	response, getErr := g.client.get(urlStr, &result, nil)

	return &result, response, getErr
}

// Get a group by their group Id.
//
// Reference: https://docs.split.io/reference#get-group
func (g *GroupsService) Get(id string) (*Group, *simpleresty.Response, error) {
	var result Group
	urlStr := g.client.http.RequestURL("/groups/%s", id)

	response, getErr := g.client.get(urlStr, &result, nil)

	return &result, response, getErr
}

// Create a new group in your organization.
//
// Reference: https://docs.split.io/reference#create-group
func (g *GroupsService) Create(opts *GroupRequest) (*Group, *simpleresty.Response, error) {
	var result Group
	urlStr := g.client.http.RequestURL("/groups")

	// Execute the request
	response, getErr := g.client.post(urlStr, &result, opts)

	return &result, response, getErr
}

// Update a group in your organization.
//
// Reference: https://docs.split.io/reference#update-group
func (g *GroupsService) Update(id string, opts *GroupRequest) (*Group, *simpleresty.Response, error) {
	var result Group
	urlStr := g.client.http.RequestURL("/groups/%s", id)

	// Execute the request
	response, getErr := g.client.put(urlStr, &result, opts)

	return &result, response, getErr
}

// Delete a group in your organization.
//
// Reference: https://docs.split.io/reference#delete-group
func (g *GroupsService) Delete(id string) (*simpleresty.Response, error) {
	urlStr := g.client.http.RequestURL("/groups/%s", id)

	// Execute the request
	response, getErr := g.client.delete(urlStr, nil, nil)

	return response, getErr
}
