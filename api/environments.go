package api

import (
	"fmt"
	"github.com/davidji99/simpleresty"
	"strconv"
)

// EnvironmentsService handles communication with the environments related
// methods of the Split.io APIv2.
//
// Reference: https://docs.split.io/reference#environments-overview
type EnvironmentsService service

// Environment reflects a stage in the development process, such as your production application or your internal staging
// environment. During the feature release process, Splits can be promoted through the various environments; allowing for
// a targeted roll out throughout the development process.
type Environment struct {
	ID         *string `json:"id"`
	Name       *string `json:"name"`
	Production *bool   `json:"production"`
}

// EnvironmentRequest represents a request modify an environment.
type EnvironmentRequest struct {
	Name       *string `json:"name,omitempty"`
	Production *bool   `json:"production,omitempty"`
}

// List all environments.
//
// Reference: https://docs.split.io/reference#get-environments
func (e *EnvironmentsService) List(workspaceID string) ([]*Environment, *simpleresty.Response, error) {
	var result []*Environment
	urlStr := e.client.http.RequestURL("/environments/ws/%s", workspaceID)

	// Execute the request
	response, getErr := e.client.http.Get(urlStr, &result, nil)

	return result, response, getErr
}

// ListSegments retrieves the Segments given an environment.
//
// Reference: https://docs.split.io/reference#list-segments-in-environment
func (e *EnvironmentsService) ListSegments(workspaceID string) (*SegmentListResult, *simpleresty.Response, error) {
	var result SegmentListResult
	urlStr := e.client.http.RequestURL("/segments/ws/%s", workspaceID)

	// Execute the request
	response, getErr := e.client.http.Get(urlStr, &result, nil)

	return &result, response, getErr
}

// FindByID retrieves an environment by its ID.
//
// Note: this method uses the List() method to first return all environments and then look for the target environment
// by an ID. The Split APIv2 does not provide a GET#show endpoint for environments unfortunately.
func (e *EnvironmentsService) FindByID(workspaceID, envID string) (*Environment, *simpleresty.Response, error) {
	envs, listResponse, listErr := e.List(workspaceID)
	if listErr != nil {
		return nil, listResponse, listErr
	}

	for _, e := range envs {
		if e.GetID() == envID {
			return e, nil, nil
		}
	}

	return nil, nil, fmt.Errorf("environment [%s] not found", envID)
}

// Create an environment.
//
// Reference: https://docs.split.io/reference#create-environment
func (e *EnvironmentsService) Create(workspaceID string, opts *EnvironmentRequest) (*Environment, *simpleresty.Response, error) {
	var result Environment
	urlStr := e.client.http.RequestURL("/environments/ws/%s", workspaceID)

	// Execute the request
	response, createErr := e.client.http.Post(urlStr, &result, opts)

	return &result, response, createErr
}

type environmentPatchRequest struct {
	Operation string `json:"op"`
	Path      string `json:"path"`
	Value     string `json:"value"`
}

// Update an environment.
//
// Reference: https://docs.split.io/reference#update-environment
func (e *EnvironmentsService) Update(workspaceID, envID string, opts *EnvironmentRequest) (*Environment, *simpleresty.Response, error) {
	var result Environment
	urlStr := e.client.http.RequestURL("/environments/ws/%s/%s", workspaceID, envID)

	// Construct request body.
	reqBody := make([]environmentPatchRequest, 0)

	if opts.Name != nil {
		reqBody = append(reqBody, environmentPatchRequest{
			Operation: "replace",
			Path:      "/name",
			Value:     opts.GetName(),
		})
	}

	if opts.Production != nil {
		reqBody = append(reqBody, environmentPatchRequest{
			Operation: "replace",
			Path:      "/production",
			Value:     strconv.FormatBool(opts.GetProduction()),
		})
	}

	// Execute the request
	response, getErr := e.client.http.Patch(urlStr, &result, reqBody)

	return &result, response, getErr
}

// Delete an environment.
//
// Note: you CANNOT delete an environment unless you first revoke all api keys associated with it.
//
// If deletion request is successful, the response body returns a "true" string.
//
// Reference: https://docs.split.io/reference#delete-environment
func (e *EnvironmentsService) Delete(workspaceID, envID string) (*simpleresty.Response, error) {
	urlStr := e.client.http.RequestURL("/environments/ws/%s/%s", workspaceID, envID)
	// Execute the request
	response, getErr := e.client.http.Delete(urlStr, nil, nil)

	return response, getErr
}
