package api

import (
	"fmt"
	"github.com/davidji99/simpleresty"
)

// WorkspacesService handles communication with the workspace related
// methods of the Split.io APIv2.
//
// Reference: https://docs.split.io/reference/workspaces-overview
type WorkspacesService service

// Workspaces represents all workspaces.
type Workspaces struct {
	Objects []*Workspace `json:"objects"`
	GenericListResult
}

// Workspace represents a workspace.
type Workspace struct {
	Name                     *string `json:"name"`
	Type                     *string `json:"type"`
	ID                       *string `json:"id"`
	RequiresTitleAndComments *bool   `json:"requiresTitleAndComments"`
}

// WorkspaceRequest represents a request to create/update a workspace.
type WorkspaceRequest struct {
	Name                     *string `json:"name,omitempty"`
	RequiresTitleAndComments *bool   `json:"requiresTitleAndComments,omitempty"` // Require title and comments for splits, segment, and metric changes.
}

// workspaceUpdateRequest represents the full request to update an existing workspace
type workspaceUpdateRequestFull struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

// WorkspaceListQueryParams represents query parameters when retrieving all workspaces.
type WorkspaceListQueryParams struct {
	// The offset to retrieve. Useful for pagination
	Offset int `url:"offset,omitempty"`

	// The maximum limit to return per call. Max=20.
	Limit int `url:"limit,omitempty"`

	// Filter workspaces by name.
	Name string `url:"name,omitempty"`

	// Match operator for the name query parameter. `IS` is the default value but STARTS_WITH and CONTAINS is also supported.
	NameOp string `url:"nameOp,omitempty"`
}

// List all workspaces.
//
// Reference: https://docs.split.io/reference#get-workspaces
func (w *WorkspacesService) List(opts ...interface{}) (*Workspaces, *simpleresty.Response, error) {
	var result *Workspaces
	urlStr, err := w.client.http.RequestURLWithQueryParams("/workspaces", opts...)
	if err != nil {
		return nil, nil, err
	}

	// Execute the request
	response, getErr := w.client.http.Get(urlStr, &result, nil)

	return result, response, getErr
}

// FindById retrieves a workspace by its ID.
//
// Note: this method uses the List() method to first return all workspaces and then look for the target workspace
// by an ID. The Split APIv2 does not provide a GET#show endpoint for workspaces, unfortunately.
func (w *WorkspacesService) FindById(id string) (*Workspace, *simpleresty.Response, error) {
	workspaces, listResponse, listErr := w.List()
	if listErr != nil {
		return nil, listResponse, listErr
	}

	if workspaces.HasObjects() {
		for _, w := range workspaces.Objects {
			if w.GetID() == id {
				return w, nil, nil
			}
		}
	}

	return nil, nil, fmt.Errorf("workspace [%s] not found", id)
}

// FindByName retrieves a workspace by its name.
//
// This method uses List query parameters to find an exact match.
func (w *WorkspacesService) FindByName(name string) (*Workspace, *simpleresty.Response, error) {
	params := WorkspaceListQueryParams{Name: name, NameOp: "IS"}

	workspaces, listResponse, listErr := w.List(params)
	if listErr != nil {
		return nil, listResponse, listErr
	}

	if workspaces.HasObjects() {
		if len(workspaces.Objects) == 1 {
			return workspaces.Objects[0], nil, nil
		}
	}

	return nil, nil, fmt.Errorf("workspace [%s] not found", name)
}

// Create a workspaces.
//
// Note: When you create a workspace from this API, this won't create the default environment.
// You must use the create environment API to create an environment.
//
// Reference: https://docs.split.io/reference/create-workspace
func (w *WorkspacesService) Create(opts *WorkspaceRequest) (*Workspace, *simpleresty.Response, error) {
	var result Workspace
	urlStr := w.client.http.RequestURL("/workspaces")

	// Execute the request
	response, createErr := w.client.http.Post(urlStr, &result, opts)

	return &result, response, createErr
}

// Update a workspaces.
//
// Reference: https://docs.split.io/reference/update-workspace
func (w *WorkspacesService) Update(id string, opts *WorkspaceRequest) (*Workspace, *simpleresty.Response, error) {
	var result Workspace
	urlStr := w.client.http.RequestURL("/workspaces/%s", id)

	optsFull := make([]workspaceUpdateRequestFull, 0)
	if opts.Name != nil {
		optsFull = append(optsFull, workspaceUpdateRequestFull{
			Op:    "replace",
			Path:  "/name",
			Value: *opts.Name,
		})
	}

	if opts.RequiresTitleAndComments != nil {
		optsFull = append(optsFull, workspaceUpdateRequestFull{
			Op:    "replace",
			Path:  "/requiresTitleAndComments",
			Value: *opts.RequiresTitleAndComments,
		})
	}

	// Execute the request
	response, updateErr := w.client.http.Patch(urlStr, &result, optsFull)

	return &result, response, updateErr
}

// Delete a workspace.
//
// Reference: https://docs.split.io/reference/delete-workspace
func (w *WorkspacesService) Delete(id string) (*simpleresty.Response, error) {
	urlStr := w.client.http.RequestURL("/workspaces/%s", id)

	// Execute the request
	response, deleteErr := w.client.http.Delete(urlStr, nil, nil)

	return response, deleteErr
}
