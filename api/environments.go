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

// EnvironmentSegment represents a segment in an environment.
type EnvironmentSegment struct {
	ID             *string      `json:"id"`
	OrgID          *string      `json:"orgId"`
	Environment    *string      `json:"environment"`
	Name           *string      `json:"name"`
	TrafficTypeID  *string      `json:"trafficTypeId"`
	Description    *string      `json:"description"`
	Status         *string      `json:"status"`
	CreationTime   *int64       `json:"creationTime"`
	LastUpdateTime *int64       `json:"lastUpdateTime"`
	TrafficTypeURN *TrafficType `json:"trafficTypeURN"`
	Creator        *User        `json:"creator"`
}

// EnvironmentRequest represents a request modify an environment.
type EnvironmentRequest struct {
	Name       *string `json:"name,omitempty"`
	Production *bool   `json:"production,omitempty"`
}

// EnvironmentSegmentKeysRequest represents a request to add/remove segment keys in an environment.
type EnvironmentSegmentKeysRequest struct {
	Keys    []string `json:"keys"`
	Comment string   `json:"comment,omitempty"`
}

// List all environments.
//
// Reference: https://docs.split.io/reference#get-environments
func (e *EnvironmentsService) List(workspaceID string) ([]*Environment, *simpleresty.Response, error) {
	var result []*Environment
	urlStr := e.client.http.RequestURL("/environments/ws/%s", workspaceID)
	response, getErr := e.client.get(urlStr, &result, nil)

	return result, response, getErr
}

// ListSegments retrieves segments given an environment.
//
// Reference: https://docs.split.io/reference/list-segments-in-environment
func (e *EnvironmentsService) ListSegments(workspaceID, environmentID string) (*SegmentListResult, *simpleresty.Response, error) {
	var result SegmentListResult
	urlStr := e.client.http.RequestURL("/segments/ws/%s/environments/%s", workspaceID, environmentID)
	response, getErr := e.client.get(urlStr, &result, nil)

	return &result, response, getErr
}

// AddSegmentKeys for a given an environment.
//
// Reference: https://docs.split.io/reference/update-segment-keys-in-environment-via-json
func (e *EnvironmentsService) AddSegmentKeys(environmentID, segmentName string, shouldReplace bool, opts *EnvironmentSegmentKeysRequest) (*EnvironmentSegment, *simpleresty.Response, error) {
	var result EnvironmentSegment
	urlStr := e.client.http.RequestURL("/segments/%s/%s/uploadKeys?replace=%v", environmentID, segmentName, shouldReplace)
	response, updateErr := e.client.put(urlStr, &result, opts)

	return &result, response, updateErr
}

// GetSegmentKeys retrieves segment keys given an environment.
//
// Reference: https://docs.split.io/reference/get-segment-keys-in-environment
func (e *EnvironmentsService) GetSegmentKeys(environmentID, segmentName string) (*SegmentKeysList, *simpleresty.Response, error) {
	var result SegmentKeysList
	urlStr := e.client.http.RequestURL("/segments/%s/%s/keys", environmentID, segmentName)
	response, getErr := e.client.get(urlStr, &result, nil)

	return &result, response, getErr
}

// RemoveSegmentKeys removes segment keys given an environment.
//
// Reference: https://docs.split.io/reference/remove-segment-keys-from-environment
func (e *EnvironmentsService) RemoveSegmentKeys(environmentID, segmentName string, opts *EnvironmentSegmentKeysRequest) (*simpleresty.Response, error) {
	urlStr := e.client.http.RequestURL("/segments/%s/%s/removeKeys", environmentID, segmentName)
	response, err := e.client.put(urlStr, nil, opts)

	return response, err
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

// FindByName retrieves an environment by its name.
//
// Note: this method uses the List() method to first return all environments and then look for the target environment
// by an name. The Split APIv2 does not provide a GET#show endpoint for environments unfortunately.
func (e *EnvironmentsService) FindByName(workspaceID, envName string) (*Environment, *simpleresty.Response, error) {
	envs, listResponse, listErr := e.List(workspaceID)
	if listErr != nil {
		return nil, listResponse, listErr
	}

	for _, e := range envs {
		if e.GetName() == envName {
			return e, nil, nil
		}
	}

	return nil, nil, fmt.Errorf("environment [%s] not found", envName)
}

// Create an environment.
//
// Reference: https://docs.split.io/reference#create-environment
func (e *EnvironmentsService) Create(workspaceID string, opts *EnvironmentRequest) (*Environment, *simpleresty.Response, error) {
	var result Environment
	urlStr := e.client.http.RequestURL("/environments/ws/%s", workspaceID)

	// Execute the request
	response, createErr := e.client.post(urlStr, &result, opts)

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
	response, getErr := e.client.patch(urlStr, &result, reqBody)

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
	response, getErr := e.client.delete(urlStr, nil, nil)

	return response, getErr
}
