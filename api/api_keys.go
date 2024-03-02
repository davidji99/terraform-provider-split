package api

import "github.com/davidji99/simpleresty"

var (
	ValidApiKeyTypes = []string{"client_side", "server_side", "admin"}
	ValidApiKeyRoles = []string{"API_ALL_GRANTED", "API_APIKEY", "API_ADMIN", "API_WORKSPACE_ADMIN", "API_FEATURE_FLAG_VIEWER",
		"API_FEATURE_FLAG_EDITOR", "API_SEGMENT_VIEWER", "API_SEGMENT_EDITOR"}
)

// KeysService handles communication with the api keys related
// methods of the Split.io APIv2.
//
// Reference: https://docs.split.io/reference/api-keys-overview
type KeysService service

// KeyRequest represents a request to create an API key.
type KeyRequest struct {
	Name         string                  `json:"name"`
	KeyType      string                  `json:"apiKeyType"`
	Roles        []string                `json:"roles"`
	Environments []KeyEnvironmentRequest `json:"environments"`
	Workspace    *KeyWorkspaceRequest    `json:"workspace"`
}

type KeyEnvironmentRequest struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

type KeyWorkspaceRequest struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

// KeyResponse represents the created key.
//
// Not all fields are added here.
type KeyResponse struct {
	Id         *string  `json:"id"`
	Name       *string  `json:"name"`
	Roles      []string `json:"roles"`
	Type       *string  `json:"type"`
	ApiKeyType *string  `json:"apiKeyType"`

	// Key is the actual API key
	Key *string `json:"key"`
}

// Create an API key.
//
// Reference: https://docs.split.io/reference/create-an-api-key
func (k *KeysService) Create(opts *KeyRequest) (*KeyResponse, *simpleresty.Response, error) {
	var result KeyResponse
	urlStr := k.client.http.RequestURL("/apiKeys")

	// Execute the request
	response, createErr := k.client.post(urlStr, &result, opts)

	return &result, response, createErr
}

// Delete an API key.
//
// Reference: https://docs.split.io/reference/delete-an-api-key
func (k *KeysService) Delete(key string) (*simpleresty.Response, error) {
	urlStr := k.client.http.RequestURL("/apiKeys/%s", key)
	// Execute the request
	response, err := k.client.delete(urlStr, nil, nil)

	return response, err
}
