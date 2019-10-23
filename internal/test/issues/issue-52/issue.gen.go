// Package issue_52 provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package issue_52

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ArrayValue defines model for ArrayValue.
type ArrayValue []Value

// Document defines model for Document.
type Document struct {
	Fields *Document_Fields `json:"fields,omitempty"`
}

// Document_Fields defines model for Document.Fields.
type Document_Fields struct {
	AdditionalProperties map[string]Value `json:"-"`
}

// Value defines model for Value.
type Value struct {
	ArrayValue  *ArrayValue `json:"arrayValue,omitempty"`
	StringValue *string     `json:"stringValue,omitempty"`
}

// Getter for additional properties for Document_Fields. Returns the specified
// element and whether it was found
func (a Document_Fields) Get(fieldName string) (value Value, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for Document_Fields
func (a *Document_Fields) Set(fieldName string, value Value) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]Value)
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for Document_Fields to handle AdditionalProperties
func (a *Document_Fields) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]Value)
		for fieldName, fieldBuf := range object {
			var fieldVal Value
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("error unmarshaling field %s", fieldName))
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for Document_Fields to handle AdditionalProperties
func (a Document_Fields) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error marshaling '%s'", fieldName))
		}
	}
	return json.Marshal(object)
}

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(req *http.Request, ctx context.Context) error

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example.
	Server string

	// HTTP client with any customized settings, such as certificate chains.
	Client *http.Client

	// A callback for modifying requests which are generated before sending over
	// the network.
	RequestEditor RequestEditorFn

	// userAgent to use
	userAgent string

	// timeout of single request
	requestTimeout time.Duration

	// timeout of idle http connections
	idleTimeout time.Duration

	// maxium idle connections of the underlying http-client.
	maxIdleConns int
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// The interface specification for the client above.
type ClientInterface interface {
	// ExampleGet request
	ExampleGet(ctx context.Context) (*http.Response, error)
}

func (c *Client) ExampleGet(ctx context.Context) (*http.Response, error) {
	req, err := NewExampleGetRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(req, ctx)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

// NewExampleGetRequest generates requests for ExampleGet
func NewExampleGetRequest(server string) (*http.Request, error) {
	var err error

	queryUrl := fmt.Sprintf("%s/example", server)

	req, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClient creates a new Client.
func NewClient(ctx context.Context, opts ...ClientOption) (*ClientWithResponses, error) {
	// create a client with sane default values
	client := Client{
		// must have a slash in order to resolve relative paths correctly.
		Server:         "",
		userAgent:      "oapi-codegen",
		maxIdleConns:   10,
		requestTimeout: 5 * time.Second,
		idleTimeout:    30 * time.Second,
	}
	// mutate defaultClient and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}

	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = client.newHTTPClient()
	}

	return &ClientWithResponses{
		ClientInterface: &client,
	}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		if !strings.HasSuffix(baseURL, "/") {
			baseURL += "/"
		}
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// WithUserAgent allows setting the userAgent
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) error {
		c.userAgent = userAgent
		return nil
	}
}

// WithIdleTimeout overrides the timeout of idle connections.
func WithIdleTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) error {
		c.idleTimeout = timeout
		return nil
	}
}

// WithRequestTimeout overrides the timeout of individual requests.
func WithRequestTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) error {
		c.requestTimeout = timeout
		return nil
	}
}

// WithMaxIdleConnections overrides the amount of idle connections of the
// underlying http-client.
func WithMaxIdleConnections(maxIdleConns uint) ClientOption {
	return func(c *Client) error {
		c.maxIdleConns = int(maxIdleConns)
		return nil
	}
}

// WithHTTPClient allows overriding the default httpClient, which is
// automatically created. This is useful for tests.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) error {
		c.Client = httpClient
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditor = fn
		return nil
	}
}

// newHTTPClient creates a httpClient for the current connection options.
func (c *Client) newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: c.requestTimeout,
		Transport: &http.Transport{
			MaxIdleConns:    c.maxIdleConns,
			IdleConnTimeout: c.idleTimeout,
		},
	}
}

// NewClientWithResponses returns a ClientWithResponses with a default Client:
func NewClientWithResponses(server string) *ClientWithResponses {
	return &ClientWithResponses{
		ClientInterface: &Client{
			Client: &http.Client{},
			Server: server,
		},
	}
}

// NewClientWithResponsesAndRequestEditorFunc takes in a RequestEditorFn callback function and returns a ClientWithResponses with a default Client:
func NewClientWithResponsesAndRequestEditorFunc(server string, reqEditorFn RequestEditorFn) *ClientWithResponses {
	return &ClientWithResponses{
		ClientInterface: &Client{
			Client:        &http.Client{},
			Server:        server,
			RequestEditor: reqEditorFn,
		},
	}
}

type exampleGetResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Document
}

// Status returns HTTPResponse.Status
func (r exampleGetResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r exampleGetResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// ExampleGetWithResponse request returning *ExampleGetResponse
func (c *ClientWithResponses) ExampleGetWithResponse(ctx context.Context) (*exampleGetResponse, error) {
	rsp, err := c.ExampleGet(ctx)
	if err != nil {
		return nil, err
	}
	return ParseexampleGetResponse(rsp)
}

// ParseexampleGetResponse parses an HTTP response from a ExampleGetWithResponse call
func ParseexampleGetResponse(rsp *http.Response) (*exampleGetResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &exampleGetResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		response.JSON200 = &Document{}
		if err := json.Unmarshal(bodyBytes, response.JSON200); err != nil {
			return nil, err
		}

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// (GET /example)
	ExampleGet(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// ExampleGet converts echo context to params.
func (w *ServerInterfaceWrapper) ExampleGet(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ExampleGet(ctx)
	return err
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router runtime.EchoRouter, si ServerInterface) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET("/example", wrapper.ExampleGet)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/5RST0v8MBD9KmV+v2No63rrTRBERPTkycuYzG6zpklIpovL0u8uk+5fFMVTkse8N29e",
	"Zgc6DDF48pyh20HWPQ1Yrjcp4fYF3UjyskxDgf8nWkIH/5oTsdmzmrl6UsDbSNABioS8b4MeB/IsAjGF",
	"SIktFbmlJWfKDY2xbINH93xR8ZeG4W1NmmH6iig4jnJpAC/G/KnZWSCTgszJ+tWRuG83o98ZEMj6ZZBi",
	"QznrZKOMCx084jtVeUxUcY9cJdJjynZDlUjkChNVPXrjyFSzd7d99aCALTtpQR84REegYEMpz5pt3dZX",
	"4jNE8hgtdHBdt/UCFETkvozeHIjdDlZUPkfUUWzdm5PwHTEoSJRj8HlObdG2cujgef+tGKOzunCbdRYP",
	"h236LdfjcpSMDJ1H8/Qg6DRNnwEAAP//F7YifqkCAAA=",
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file.
func GetSwagger() (*openapi3.Swagger, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error loading Swagger: %s", err)
	}
	return swagger, nil
}