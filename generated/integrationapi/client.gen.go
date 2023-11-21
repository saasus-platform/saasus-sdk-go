// Package integrationapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package integrationapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// ReturnInternalServerError request
	ReturnInternalServerError(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateEventBridgeEvent request with any body
	CreateEventBridgeEventWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateEventBridgeEvent(ctx context.Context, body CreateEventBridgeEventJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteEventBridgeSettings request
	DeleteEventBridgeSettings(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetEventBridgeSettings request
	GetEventBridgeSettings(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// SaveEventBridgeSettings request with any body
	SaveEventBridgeSettingsWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	SaveEventBridgeSettings(ctx context.Context, body SaveEventBridgeSettingsJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateEventBridgeTestEvent request
	CreateEventBridgeTestEvent(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) ReturnInternalServerError(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewReturnInternalServerErrorRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateEventBridgeEventWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateEventBridgeEventRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateEventBridgeEvent(ctx context.Context, body CreateEventBridgeEventJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateEventBridgeEventRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteEventBridgeSettings(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteEventBridgeSettingsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetEventBridgeSettings(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetEventBridgeSettingsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) SaveEventBridgeSettingsWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSaveEventBridgeSettingsRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) SaveEventBridgeSettings(ctx context.Context, body SaveEventBridgeSettingsJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSaveEventBridgeSettingsRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateEventBridgeTestEvent(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateEventBridgeTestEventRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewReturnInternalServerErrorRequest generates requests for ReturnInternalServerError
func NewReturnInternalServerErrorRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/errors/internal-server-error")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewCreateEventBridgeEventRequest calls the generic CreateEventBridgeEvent builder with application/json body
func NewCreateEventBridgeEventRequest(server string, body CreateEventBridgeEventJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateEventBridgeEventRequestWithBody(server, "application/json", bodyReader)
}

// NewCreateEventBridgeEventRequestWithBody generates requests for CreateEventBridgeEvent with any type of body
func NewCreateEventBridgeEventRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/eventbridge/event")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewDeleteEventBridgeSettingsRequest generates requests for DeleteEventBridgeSettings
func NewDeleteEventBridgeSettingsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/eventbridge/info")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetEventBridgeSettingsRequest generates requests for GetEventBridgeSettings
func NewGetEventBridgeSettingsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/eventbridge/info")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewSaveEventBridgeSettingsRequest calls the generic SaveEventBridgeSettings builder with application/json body
func NewSaveEventBridgeSettingsRequest(server string, body SaveEventBridgeSettingsJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewSaveEventBridgeSettingsRequestWithBody(server, "application/json", bodyReader)
}

// NewSaveEventBridgeSettingsRequestWithBody generates requests for SaveEventBridgeSettings with any type of body
func NewSaveEventBridgeSettingsRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/eventbridge/info")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewCreateEventBridgeTestEventRequest generates requests for CreateEventBridgeTestEvent
func NewCreateEventBridgeTestEventRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/eventbridge/test-event")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// ReturnInternalServerError request
	ReturnInternalServerErrorWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*ReturnInternalServerErrorResponse, error)

	// CreateEventBridgeEvent request with any body
	CreateEventBridgeEventWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateEventBridgeEventResponse, error)

	CreateEventBridgeEventWithResponse(ctx context.Context, body CreateEventBridgeEventJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateEventBridgeEventResponse, error)

	// DeleteEventBridgeSettings request
	DeleteEventBridgeSettingsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*DeleteEventBridgeSettingsResponse, error)

	// GetEventBridgeSettings request
	GetEventBridgeSettingsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetEventBridgeSettingsResponse, error)

	// SaveEventBridgeSettings request with any body
	SaveEventBridgeSettingsWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*SaveEventBridgeSettingsResponse, error)

	SaveEventBridgeSettingsWithResponse(ctx context.Context, body SaveEventBridgeSettingsJSONRequestBody, reqEditors ...RequestEditorFn) (*SaveEventBridgeSettingsResponse, error)

	// CreateEventBridgeTestEvent request
	CreateEventBridgeTestEventWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*CreateEventBridgeTestEventResponse, error)
}

type ReturnInternalServerErrorResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r ReturnInternalServerErrorResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ReturnInternalServerErrorResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateEventBridgeEventResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r CreateEventBridgeEventResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateEventBridgeEventResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteEventBridgeSettingsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r DeleteEventBridgeSettingsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteEventBridgeSettingsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetEventBridgeSettingsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *EventBridgeSettings
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetEventBridgeSettingsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetEventBridgeSettingsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type SaveEventBridgeSettingsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r SaveEventBridgeSettingsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r SaveEventBridgeSettingsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateEventBridgeTestEventResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r CreateEventBridgeTestEventResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateEventBridgeTestEventResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// ReturnInternalServerErrorWithResponse request returning *ReturnInternalServerErrorResponse
func (c *ClientWithResponses) ReturnInternalServerErrorWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*ReturnInternalServerErrorResponse, error) {
	rsp, err := c.ReturnInternalServerError(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseReturnInternalServerErrorResponse(rsp)
}

// CreateEventBridgeEventWithBodyWithResponse request with arbitrary body returning *CreateEventBridgeEventResponse
func (c *ClientWithResponses) CreateEventBridgeEventWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateEventBridgeEventResponse, error) {
	rsp, err := c.CreateEventBridgeEventWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateEventBridgeEventResponse(rsp)
}

func (c *ClientWithResponses) CreateEventBridgeEventWithResponse(ctx context.Context, body CreateEventBridgeEventJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateEventBridgeEventResponse, error) {
	rsp, err := c.CreateEventBridgeEvent(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateEventBridgeEventResponse(rsp)
}

// DeleteEventBridgeSettingsWithResponse request returning *DeleteEventBridgeSettingsResponse
func (c *ClientWithResponses) DeleteEventBridgeSettingsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*DeleteEventBridgeSettingsResponse, error) {
	rsp, err := c.DeleteEventBridgeSettings(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteEventBridgeSettingsResponse(rsp)
}

// GetEventBridgeSettingsWithResponse request returning *GetEventBridgeSettingsResponse
func (c *ClientWithResponses) GetEventBridgeSettingsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetEventBridgeSettingsResponse, error) {
	rsp, err := c.GetEventBridgeSettings(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetEventBridgeSettingsResponse(rsp)
}

// SaveEventBridgeSettingsWithBodyWithResponse request with arbitrary body returning *SaveEventBridgeSettingsResponse
func (c *ClientWithResponses) SaveEventBridgeSettingsWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*SaveEventBridgeSettingsResponse, error) {
	rsp, err := c.SaveEventBridgeSettingsWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSaveEventBridgeSettingsResponse(rsp)
}

func (c *ClientWithResponses) SaveEventBridgeSettingsWithResponse(ctx context.Context, body SaveEventBridgeSettingsJSONRequestBody, reqEditors ...RequestEditorFn) (*SaveEventBridgeSettingsResponse, error) {
	rsp, err := c.SaveEventBridgeSettings(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSaveEventBridgeSettingsResponse(rsp)
}

// CreateEventBridgeTestEventWithResponse request returning *CreateEventBridgeTestEventResponse
func (c *ClientWithResponses) CreateEventBridgeTestEventWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*CreateEventBridgeTestEventResponse, error) {
	rsp, err := c.CreateEventBridgeTestEvent(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateEventBridgeTestEventResponse(rsp)
}

// ParseReturnInternalServerErrorResponse parses an HTTP response from a ReturnInternalServerErrorWithResponse call
func ParseReturnInternalServerErrorResponse(rsp *http.Response) (*ReturnInternalServerErrorResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ReturnInternalServerErrorResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseCreateEventBridgeEventResponse parses an HTTP response from a CreateEventBridgeEventWithResponse call
func ParseCreateEventBridgeEventResponse(rsp *http.Response) (*CreateEventBridgeEventResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateEventBridgeEventResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseDeleteEventBridgeSettingsResponse parses an HTTP response from a DeleteEventBridgeSettingsWithResponse call
func ParseDeleteEventBridgeSettingsResponse(rsp *http.Response) (*DeleteEventBridgeSettingsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteEventBridgeSettingsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseGetEventBridgeSettingsResponse parses an HTTP response from a GetEventBridgeSettingsWithResponse call
func ParseGetEventBridgeSettingsResponse(rsp *http.Response) (*GetEventBridgeSettingsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetEventBridgeSettingsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest EventBridgeSettings
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseSaveEventBridgeSettingsResponse parses an HTTP response from a SaveEventBridgeSettingsWithResponse call
func ParseSaveEventBridgeSettingsResponse(rsp *http.Response) (*SaveEventBridgeSettingsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &SaveEventBridgeSettingsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseCreateEventBridgeTestEventResponse parses an HTTP response from a CreateEventBridgeTestEventWithResponse call
func ParseCreateEventBridgeTestEventResponse(rsp *http.Response) (*CreateEventBridgeTestEventResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateEventBridgeTestEventResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}
