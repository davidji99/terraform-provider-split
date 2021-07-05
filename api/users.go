package api

import "github.com/davidji99/simpleresty"

// UsersService handles communication with the user related
// methods of the Split.io APIv2.
//
// Reference: https://docs.split.io/reference#users-overview
type UsersService service

// User represents a user.
type User struct {
	ID     *string  `json:"id"`
	Type   *string  `json:"type"`
	Name   *string  `json:"name"`
	Email  *string  `json:"email"`
	Status *string  `json:"status"`
	TFA    *bool    `json:"2fa"`
	Groups []*Group `json:"groups"`
}

// UserListResult
type UserListResult struct {
	Data           []*User `json:"data"`
	NextMarker     *string `json:"nextMarker,omitempty"`
	PreviousMarker *string `json:"previousMarker,omitempty"`
	Limit          *int    `json:"limit"`
	Count          *int    `json:"count"`
}

// UserListOpts
type UserListOpts struct {
	// ACTIVE | DEACTIVATED | PENDING are the allowed status values to filter by
	Status string `url:"status,omitempty"`

	// 1-200 are the potential values. Default=50
	Limit int `url:"limit,omitempty"`

	// value of "previousMarker" in response
	Before int `url:"limit,omitempty"`

	// value of "nextMarker" in response
	After string `url:"limit,omitempty"`

	// eturns Active members of a group
	GroupID string `url:"limit,omitempty"`
}

// UserCreateRequest
type UserCreateRequest struct {
	Email  string `json:"email,omitempty"`
	Groups []struct {
		ID   string `json:"id,omitempty"`
		Type string `json:"type,omitempty"`
	} `json:"groups,omitempty"`
}

type UserUpdateRequest struct {
	Name   string `json:"name,omitempty"`
	Email  string `json:"email,omitempty"`
	TFA    *bool  `json:"2fa,omitempty"`
	Status string `json:"status,omitempty"`
}

// List all active, deactivated, and pending users in the organization.
//
// Reference: https://docs.split.io/reference#list-users
func (u *UsersService) List(opts *UserListOpts) (*UserListResult, *simpleresty.Response, error) {
	var result UserListResult
	urlStr, urlStrErr := u.client.http.RequestURLWithQueryParams("/users", opts)
	if urlStrErr != nil {
		return nil, nil, urlStrErr
	}

	// Execute the request
	response, getErr := u.client.http.Get(urlStr, &result, nil)

	return &result, response, getErr
}

// Get a user by their user Id.
//
// Reference: https://docs.split.io/reference#get-user
func (u *UsersService) Get(id string) (*User, *simpleresty.Response, error) {
	var result User
	urlStr := u.client.http.RequestURL("/users/%s", id)

	// Execute the request
	response, getErr := u.client.http.Get(urlStr, &result, nil)

	return &result, response, getErr
}

// Invite a new user to your organization. They will be created with a Pending status
//
// Reference: https://docs.split.io/reference#invite-a-new-user
func (u *UsersService) Invite(opts *UserCreateRequest) (*User, *simpleresty.Response, error) {
	var result User
	urlStr := u.client.http.RequestURL("/users")

	// Execute the request
	response, err := u.client.http.Post(urlStr, &result, opts)

	return &result, response, err
}

// Update display name, email, disable 2FA, and Activate/Deactivate of a User.
//
// Reference: https://docs.split.io/reference#full-update-user
func (u *UsersService) Update(id string, opts *UserUpdateRequest) (*User, *simpleresty.Response, error) {
	var result User
	urlStr := u.client.http.RequestURL("/users/%s", id)

	// Execute the request
	response, err := u.client.http.Put(urlStr, &result, opts)

	return &result, response, err
}

//// UpdateUserGroups Use this endpoint to update the groups that a user is part of.
////
//// Reference: https://docs.split.io/reference#update-users-groups
//func (s *UsersService) UpdateUserGroups(id string, opts *UserUpdateRequest) (*User, *simpleresty.Response, error) {
//	var result User
//	urlStr := s.client.http.RequestURL("/users/%s", id)
//
//	// Execute the request
//	response, err := s.client.http.Put(urlStr, &result, opts)
//
//	return &result, response, err
//}

// DeletePendingUser that have not accepted their invites yet. Once a user is active,
// you can only deactivate the user via a PUT request
//
// Reference: https://docs.split.io/reference#delete-a-pending-user
func (u *UsersService) DeletePendingUser(id string) (*simpleresty.Response, error) {
	urlStr := u.client.http.RequestURL("/users/%s", id)

	// Execute the request
	response, deleteErr := u.client.http.Delete(urlStr, nil, nil)

	return response, deleteErr
}
