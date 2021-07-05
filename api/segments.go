package api

import "github.com/davidji99/simpleresty"

// SegmentsService handles communication with the segment related
// methods of the Split.io APIv2.
//
// Reference: https://docs.split.io/reference#segments-overview
type SegmentsService service

// Segment represents a Split segment.
type Segment struct {
	Name         *string      `json:"name"`
	Description  *string      `json:"description,omitempty"`
	Environment  *Environment `json:"environment"`
	TrafficType  *TrafficType `json:"trafficType"`
	CreationTime *int64       `json:"creationTime"`
	Tags         []*Tag       `json:"tags"`
}

// SegmentListResult represents the response returned when listing all segments.
type SegmentListResult struct {
	Objects []*Segment `json:"objects"`
	GenericListResult
}

// SegmentRequest represents a request to create a segment.
type SegmentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// List all segments.
//
// Reference: https://docs.split.io/reference#list-segments
func (s *SegmentsService) List(workspaceID string) (*SegmentListResult, *simpleresty.Response, error) {
	var result SegmentListResult
	urlStr := s.client.http.RequestURL("/segments/ws/%s", workspaceID)

	// Execute the request
	response, getErr := s.client.http.Get(urlStr, &result, nil)

	return &result, response, getErr
}

// Get a segment.
//
// Reference: n/a
func (s *SegmentsService) Get(workspaceID, name string) (*Segment, *simpleresty.Response, error) {
	var result Segment
	urlStr := s.client.http.RequestURL("/segments/ws/%s/%s", workspaceID, name)

	// Execute the request
	response, getErr := s.client.http.Get(urlStr, &result, nil)

	return &result, response, getErr
}

// Create a segment. This API does not configure the Segment in any environment.
//
// Reference: https://docs.split.io/reference#create-segment
func (s *SegmentsService) Create(workspaceID, trafficTypeID string, opts *SegmentRequest) (*Segment, *simpleresty.Response, error) {
	var result Segment
	urlStr := s.client.http.RequestURL("/segments/ws/%s/trafficTypes/%s", workspaceID, trafficTypeID)

	// Execute the request
	response, createErr := s.client.http.Post(urlStr, &result, opts)

	return &result, response, createErr
}

// There exists an update (PUT) endpoint to modify the description but that is an internal endpoint:
// https://app.split.io/internal/api/segmentMetadata/updateDescription/<SEGMENT_ID>
//
// TODO: implement this method whenever this endpoint is GA.
//func (s *SegmentsService) Update() {
//}

// Delete a segment. This will automatically unconfigure the Segment Definition from all environments.
//
// Reference: https://docs.split.io/reference#delete-segment
func (s *SegmentsService) Delete(workspaceID, segmentName string) (*simpleresty.Response, error) {
	urlStr := s.client.http.RequestURL("/segments/ws/%s/%s", workspaceID, segmentName)

	// Execute the request
	response, deleteErr := s.client.http.Delete(urlStr, nil, nil)

	return response, deleteErr
}

// Activate a Segment in an environment to be able to set its definitions.
//
// Reference: https://docs.split.io/reference#enable-segment-in-environment
func (s *SegmentsService) Activate(environmentID, segmentName string) (*Segment, *simpleresty.Response, error) {
	var result Segment
	urlStr := s.client.http.RequestURL("/segments/%s/%s", environmentID, segmentName)

	// Execute the request
	response, err := s.client.http.Post(urlStr, &result, nil)

	return &result, response, err
}

// Deactivate a Segment in an environment.
//
// Reference: https://docs.split.io/reference#deactivate-segment-in-environment
func (s *SegmentsService) Deactivate(environmentID, segmentName string) (*Segment, *simpleresty.Response, error) {
	var result Segment
	urlStr := s.client.http.RequestURL("/segments/%s/%s", environmentID, segmentName)

	// Execute the request
	response, err := s.client.http.Delete(urlStr, &result, nil)

	return &result, response, err
}
