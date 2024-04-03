package api

import (
	"fmt"
	"github.com/davidji99/simpleresty"
)

// FeatureFlagsService handles communication with the environments related
// methods of the Split.io APIv2.
//
// Reference: https://docs.split.io/reference/feature-flag-overview
type FeatureFlagsService service

// FeatureFlag represents a feature flag in Split.
type FeatureFlag struct {
	ID           *string      `json:"id"`
	Name         *string      `json:"name"`
	Description  *string      `json:"description"`
	TrafficType  *TrafficType `json:"trafficType"`
	CreationTime *int64       `json:"creationTime"`
	Tags         []*Tag       `json:"tags"`
}

// FeatureFlagsListResponse represents the response when listing feature flags.
type FeatureFlagsListResponse struct {
	GenericListResult
	FeatureFlags []*FeatureFlag `json:"objects"`
}

type FeatureFlagsListOpts struct {
	GenericListQueryParams

	// Tags are repeatable tag parameter(s) to query by
	Tags []string `url:"tag,omitempty"`
}

// List feature flags.
//
// Reference: https://docs.split.io/reference/list-feature-flags
func (f *FeatureFlagsService) List(workspaceID string, opts *FeatureFlagsListOpts) (*FeatureFlagsListResponse, *simpleresty.Response, error) {
	var result FeatureFlagsListResponse
	urlStr, err := f.client.http.RequestURLWithQueryParams(fmt.Sprintf("/splits/ws/%s", workspaceID), opts)
	if err != nil {
		return nil, nil, err
	}

	response, getErr := f.client.get(urlStr, &result, nil)

	return &result, response, getErr
}

// ListAll feature flags. Compared to `List()`, `ListAll()` will return all feature flags by listing through all pages.
//
// Reference: https://docs.split.io/reference/list-feature-flags
func (f *FeatureFlagsService) ListAll(workspaceID string, opts *FeatureFlagsListOpts) ([]*FeatureFlag, *simpleresty.Response, error) {
	var result FeatureFlagsListResponse
	urlStr, err := f.client.http.RequestURLWithQueryParams(fmt.Sprintf("/splits/ws/%s", workspaceID), opts)
	if err != nil {
		return nil, nil, err
	}

	response, getErr := f.client.get(urlStr, &result, nil)

	return &result, response, getErr
}
