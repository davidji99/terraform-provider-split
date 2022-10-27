package api

import (
	"fmt"

	"github.com/davidji99/simpleresty"
)

// SplitDefinitions represents all definitions in an environment.
type SplitDefinitions struct {
	GenericListResult
	Objects []*SplitDefinition `json:"objects"`
}

// SplitDefinition is the configuration of a Split in a specific environment such as production or staging.
type SplitDefinition struct {
	ID                *string      `json:"id"`
	Name              *string      `json:"name"`
	Environment       *Environment `json:"environment"`
	TrafficType       *TrafficType `json:"trafficType"`
	Killed            *bool        `json:"killed"`
	Treatments        []*Treatment `json:"treatments"`
	DefaultTreatment  *string      `json:"defaultTreatment"`
	TrafficAllocation *int         `json:"trafficAllocation"`
	Rules             []*Rule      `json:"rules"`
	DefaultRule       []*Bucket    `json:"defaultRule"`
	CreationTIme      *int         `json:"creationTIme"`
	LastUpdateTime    *int         `json:"lastUpdateTime"`
}

// SplitDefinitionRequest creates or updates a split definition.
type SplitDefinitionRequest struct {
	Treatments       []Treatment `json:"treatments"`
	DefaultTreatment string      `json:"defaultTreatment"`
	Rules            []Rule      `json:"rules"`
	DefaultRule      []Bucket    `json:"defaultRule"`
	Comment          string      `json:"comment,omitempty"`
}

// Rule consists of a Condition and a list of Buckets.
//
// When the Split Definition is evaluated, if the Condition of this Rule is met, then the customer will be evaluated
// and fall in the sticky distribution specified in the Buckets
type Rule struct {
	Condition *Condition `json:"condition"`
	Buckets   []*Bucket  `json:"buckets"`
}

// Treatment is a state of a Split. A simple feature flag has two treatments: on and off. An experiment can have any number of treatments.
type Treatment struct {
	Name           *string  `json:"name"`
	Configurations *string  `json:"configurations"`
	Description    *string  `json:"description"`
	Keys           *string  `json:"keys,omitempty"`
	Segment        []string `json:"segment,omitempty"`
}

// Bucket represents a sticky distribution of customers into treatments of a Split.
type Bucket struct {
	Treatment *string `json:"treatment"`
	Size      *int    `json:"size"`
}

// Condition consists of different Matchers. A customer satisfies a Condition only of it satisfies all the Matchers.
//
// Note that when configuring or updating a `Split Definition, there is no need to specify the combiner, since for now, Split only allows AND combiner.
type Condition struct {
	Combiner *string    `json:"combiner,omitempty"`
	Matchers []*Matcher `json:"matchers,omitempty"`
}

// Matcher represents the logic for selecting a specific subset of your customer population.
type Matcher struct {
	Negate    *bool       `json:"negate,omitempty"`
	Type      *string     `json:"type,omitempty"`
	Attribute *string     `json:"attribute,omitempty"`
	String    *string     `json:"string,omitempty"`
	Bool      *bool       `json:"bool,omitempty"`
	Strings   []string    `json:"strings,omitempty"`
	Number    *int        `json:"number,omitempty"`
	Date      *int        `json:"date,omitempty"`
	Between   interface{} `json:"between,omitempty"`
	Depends   interface{} `json:"depends,omitempty"`
}

// ListDefinitions retrieves the Split Definitions given an environment.
//
// Reference: https://docs.split.io/reference/lists-split-definitions-in-environment
func (s *SplitsService) ListDefinitions(workspaceId, environmentId string, opts ...interface{}) (*SplitDefinitions, *simpleresty.Response, error) {
	var result SplitDefinitions
	// Here I have a problem with el requestURLWITHQUERy...
	urlStr, err := s.client.http.RequestURLWithQueryParams(fmt.Sprintf("/splits/ws/%s/environments/%s", workspaceId, environmentId), opts...)
	if err != nil {
		return nil, nil, err
	}

	// Execute the request
	response, getErr := s.client.get(urlStr, &result, nil)

	return &result, response, getErr
}

// GetDefinition retrieves a Split Definition given the name and the environment.
//
// Reference: https://docs.split.io/reference/get-split-definition-in-environment
func (s *SplitsService) GetDefinition(workspaceId, splitName, environmentId string) (*SplitDefinition, *simpleresty.Response, error) {
	var result SplitDefinition
	urlStr := s.client.http.RequestURL("/splits/ws/%s/%s/environments/%s", workspaceId, splitName, environmentId)
	response, getErr := s.client.get(urlStr, &result, nil)

	return &result, response, getErr
}

// CreateDefinition configures a Split Definition for a specific environment.
//
// Reference: https://docs.split.io/reference/create-split-definition-in-environment
func (s *SplitsService) CreateDefinition(workspaceId, splitName, environmentId string, opts *SplitDefinitionRequest) (*SplitDefinition, *simpleresty.Response, error) {
	var result SplitDefinition
	urlStr := s.client.http.RequestURL("/splits/ws/%s/%s/environments/%s", workspaceId, splitName, environmentId)

	// Execute the request
	response, createErr := s.client.post(urlStr, &result, opts)

	return &result, response, createErr
}

// UpdateDefinitionFull performs a full update of a Split Definition for a specific environment.
//
// Reference: https://docs.split.io/reference/full-update-split-definition-in-environment
func (s *SplitsService) UpdateDefinitionFull(workspaceId, splitName, environmentId string, opts *SplitDefinitionRequest) (*SplitDefinition, *simpleresty.Response, error) {
	var result SplitDefinition
	urlStr := s.client.http.RequestURL("/splits/ws/%s/%s/environments/%s", workspaceId, splitName, environmentId)

	// Execute the request
	response, createErr := s.client.put(urlStr, &result, opts)

	return &result, response, createErr
}

// RemoveDefinition unconfigures a Split Definition for a specific environment.
//
// Reference: https://docs.split.io/reference/remove-split-definition-from-environment
func (s *SplitsService) RemoveDefinition(workspaceId, splitName, environmentId string) (*simpleresty.Response, error) {
	urlStr := s.client.http.RequestURL("/splits/ws/%s/%s/environments/%s", workspaceId, splitName, environmentId)

	// Execute the request
	response, createErr := s.client.delete(urlStr, nil, nil)

	return response, createErr
}
