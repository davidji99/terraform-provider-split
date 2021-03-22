package api

import (
	"fmt"
	"github.com/davidji99/simpleresty"
)

// WorkspacesService handles communication with the environments related
// methods of the Split.io APIv2.
//
// Reference: https://docs.split.io/reference#get-workspaces
type WorkspacesService service

// Workspaces represents all workspaces.
type Workspaces struct {
	Objects    []*Workspace `json:"objects"`
	Offset     *int         `json:"offset"`
	Limit      *int         `json:"limit"`
	TotalCount *int         `json:"totalCount"`
}

// Workspace represents a workspace.
type Workspace struct {
	Name                     *string `json:"name"`
	Type                     *string `json:"type"`
	ID                       *string `json:"id"`
	RequiresTitleAndComments *bool   `json:"requiresTitleAndComments"`
}

// List all workspaces.
//
// Reference: https://docs.split.io/reference#get-workspaces
func (w *WorkspacesService) List() (*Workspaces, *simpleresty.Response, error) {
	var result *Workspaces
	urlStr := w.client.http.RequestURL("/workspaces")

	// Execute the request
	response, getErr := w.client.http.Get(urlStr, &result, nil)

	return result, response, getErr
}

// FindById retrieves a workspace by its ID.
//
// Note: this method uses the List() method to first return all workspaces and then look for the target workspace
// by an ID. The Split APIv2 does not provide a GET#show endpoint for workspaces unfortunately.
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
// Note: this method uses the List() method to first return all workspaces and then look for the target workspace
// by a name.
func (w *WorkspacesService) FindByName(name string) (*Workspace, *simpleresty.Response, error) {
	workspaces, listResponse, listErr := w.List()
	if listErr != nil {
		return nil, listResponse, listErr
	}

	if workspaces.HasObjects() {
		for _, w := range workspaces.Objects {
			if w.GetName() == name {
				return w, nil, nil
			}
		}
	}

	return nil, nil, fmt.Errorf("workspace [%s] not found", name)
}
