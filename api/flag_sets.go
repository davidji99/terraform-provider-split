package api

import (
	"fmt"
	"github.com/davidji99/simpleresty"
)

// FlagSetsService handles communication with the flag sets related
// methods of the Split.io APIv3.
//
// Reference: https://docs.split.io/reference/create-flag-set
type FlagSetsService service

// FlagSet represents a flag set which allows grouping feature flags together for easier management.
type FlagSet struct {
	ID          *string         `json:"id"`
	Name        *string         `json:"name"`
	Description *string         `json:"description"`
	Workspace   *WorkspaceIDRef `json:"workspace"`
}

// WorkspaceIDRef represents a minimal representation of a workspace for flag set operations.
type WorkspaceIDRef struct {
	Type *string `json:"type"`
	ID   *string `json:"id"`
}

// FlagSetRequest represents a request to create a flag set.
type FlagSetRequest struct {
	Name        *string         `json:"name,omitempty"`
	Description *string         `json:"description,omitempty"`
	Workspace   *WorkspaceIDRef `json:"workspace,omitempty"`
}

// FlagSetListResult represents the response when listing flag sets.
type FlagSetListResult struct {
	Objects       []*FlagSet `json:"objects"`
	NextMarker    *string    `json:"nextMarker"`
	PreviousMarker *string   `json:"previousMarker"`
}

// Create a flag set.
//
// Reference: https://docs.split.io/reference/create-flag-set
func (f *FlagSetsService) Create(opts *FlagSetRequest) (*FlagSet, *simpleresty.Response, error) {
	var result FlagSet
	urlStr := "https://api.split.io/api/v3/flag-sets"
	response, createErr := f.client.post(urlStr, &result, opts)
	return &result, response, createErr
}

// FindByID retrieves a flag set by its ID.
//
// Reference: https://docs.split.io/reference/get-flag-set-by-id
func (f *FlagSetsService) FindByID(id string) (*FlagSet, *simpleresty.Response, error) {
	var result FlagSet
	urlStr := fmt.Sprintf("https://api.split.io/api/v3/flag-sets/%s", id)
	response, getErr := f.client.get(urlStr, &result, nil)
	return &result, response, getErr
}

// List all flag sets for a given workspace.
//
// Reference: https://docs.split.io/reference/list-flag-sets
func (f *FlagSetsService) List(workspaceID string) ([]*FlagSet, *simpleresty.Response, error) {
	var result FlagSetListResult
	urlStr := fmt.Sprintf("https://api.split.io/api/v3/flag-sets?workspace_id=%s", workspaceID)
	response, getErr := f.client.get(urlStr, &result, nil)
	return result.Objects, response, getErr
}

// FindByName retrieves a flag set by its name.
//
// Note: this method uses the List() method to first return all flag sets and then look for the target flag set
// by name. The Split APIv3 does not provide a direct GET endpoint to find flag sets by name.
func (f *FlagSetsService) FindByName(workspaceID, name string) (*FlagSet, *simpleresty.Response, error) {
	flagSets, listResponse, listErr := f.List(workspaceID)
	if listErr != nil {
		return nil, listResponse, listErr
	}

	for _, fs := range flagSets {
		if fs.GetName() == name {
			return fs, nil, nil
		}
	}

	return nil, nil, fmt.Errorf("flag set [%s] not found", name)
}

// Delete a flag set.
//
// Reference: https://docs.split.io/reference/delete-flag-set-by-id
func (f *FlagSetsService) Delete(id string) (*simpleresty.Response, error) {
	urlStr := fmt.Sprintf("https://api.split.io/api/v3/flag-sets/%s", id)
	response, getErr := f.client.delete(urlStr, nil, nil)
	return response, getErr
}
