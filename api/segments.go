package api

import "github.com/davidji99/simpleresty"

// SegmentsService handles communication with the segment related
// methods of the Split.io APIv2.
//
// Reference: https://docs.split.io/reference#segments-overview
type SegmentsService service

type Segment struct {
	Name         *string      `json:"name"`
	Description  *string      `json:"description"`
	Environment  *Environment `json:"environment"`
	TrafficType  *TrafficType `json:"trafficType"`
	CreationTime *int64       `json:"creationTime"`
	Tags         []*Tag       `json:"tags"`
}

type SegmentListResult struct {
	Objects []*Segment `json:"objects"`
	GenericListResult
}

type SegmentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

func (s *SegmentsService) List(workspaceID string) (*SegmentListResult, *simpleresty.Response, error) {
	var result SegmentListResult
	urlStr := s.client.http.RequestURL("/segments/ws/%s", workspaceID)

	// Execute the request
	response, getErr := s.client.http.Get(urlStr, &result, nil)

	return &result, response, getErr
}

func (s *SegmentsService) Create(workspaceID, trafficTypeID string, opts *SegmentRequest) (*Segment, *simpleresty.Response, error) {
	var result Segment
	urlStr := s.client.http.RequestURL("/segments/ws/%s/trafficTypes/%s", workspaceID, trafficTypeID)

	// Execute the request
	response, createErr := s.client.http.Post(urlStr, &result, opts)

	return &result, response, createErr
}

func (s *SegmentsService) Delete(workspaceID, segmentName string) (*simpleresty.Response, error) {
	urlStr := s.client.http.RequestURL("/segments/ws/%s/%s", workspaceID, segmentName)

	// Execute the request
	response, deleteErr := s.client.http.Delete(urlStr, nil, nil)

	return response, deleteErr
}

func (s *SegmentsService) Activate(environmentID, segmentName string) (*Segment, *simpleresty.Response, error) {
	var result Segment
	urlStr := s.client.http.RequestURL("/segments/%s/%s", environmentID, segmentName)

	// Execute the request
	response, err := s.client.http.Post(urlStr, &result, nil)

	return &result, response, err
}

func (s *SegmentsService) Deactivate(environmentID, segmentName string) (*Segment, *simpleresty.Response, error) {
	var result Segment
	urlStr := s.client.http.RequestURL("/segments/%s/%s", environmentID, segmentName)

	// Execute the request
	response, err := s.client.http.Delete(urlStr, &result, nil)

	return &result, response, err
}
