package vi

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/maxence-charriere/go-app/v11/pkg/errors"
)

// HTTPClient defines an HTTP client with configurable encoding, decoding,
// URL resolution, and header handling.
type HTTPClient interface {
	// WithEncoder sets the function used to encode request payloads.
	WithEncoder(v func(any) ([]byte, error)) HTTPClient

	// WithDecoder sets the function used to decode response payloads.
	WithDecoder(v func([]byte, any) error) HTTPClient

	// WithResolver sets the function used to resolve request URLs.
	WithResolver(v func(string) string) HTTPClient

	// WithHeader sets the function used to modify or extend HTTP headers.
	WithHeader(func(http.Header)) HTTPClient

	// Get sends an HTTP GET request to the given URL and decodes the response
	// into out.
	Get(ctx context.Context, rawURL string, out interface{}) error

	// Delete sends an HTTP DELETE request to the given URL and decodes the
	// response into out.
	Delete(ctx context.Context, rawURL string, out interface{}) error

	// Post sends an HTTP POST request with the provided input body and decodes
	// the response into out.
	Post(ctx context.Context, rawURL string, out, in interface{}) error

	// Put sends an HTTP PUT request with the provided input body and decodes
	// the response into out.
	Put(ctx context.Context, rawURL string, out, in interface{}) error

	// Do executes the given HTTP request and decodes the response into out
	// using the configured decoder. It applies any configured header and
	// URL resolver functions before sending the request. This method is
	// intended for advanced use cases that require full control over the
	// request instead of using the higher-level methods (Get, Post, Put,
	// Delete).
	Do(r *http.Request, out interface{}) error
}

// MakeHTTPClient returns a new HTTPClient using the given transport. By
// default, it uses JSON for encoding and decoding.
func MakeHTTPClient(transport http.RoundTripper) HTTPClient {
	return httpClient{
		RoundTripper: transport,
		encode:       json.Marshal,
		decode:       json.Unmarshal,
	}
}

type httpClient struct {
	http.RoundTripper

	encode    func(any) ([]byte, error)
	decode    func([]byte, any) error
	resolvers []func(string) string
	headers   []func(http.Header)
}

func (c httpClient) WithEncoder(v func(any) ([]byte, error)) HTTPClient {
	c.encode = v
	return c
}

func (c httpClient) WithDecoder(v func([]byte, any) error) HTTPClient {
	c.decode = v
	return c
}

func (c httpClient) WithResolver(v func(string) string) HTTPClient {
	c.resolvers = append(c.resolvers, v)
	return c
}

func (c httpClient) WithHeader(v func(http.Header)) HTTPClient {
	c.headers = append(c.headers, v)
	return c
}

func (c httpClient) Get(ctx context.Context, rawURL string, out interface{}) error {
	return c.get(ctx, http.MethodGet, rawURL, out)
}

func (c httpClient) Delete(ctx context.Context, rawURL string, out interface{}) error {
	return c.get(ctx, http.MethodDelete, rawURL, out)
}

func (c httpClient) Post(ctx context.Context, rawURL string, out, in interface{}) error {
	return c.post(ctx, http.MethodPost, rawURL, out, in)
}

func (c httpClient) Put(ctx context.Context, rawURL string, out, in interface{}) error {
	return c.post(ctx, http.MethodPut, rawURL, out, in)
}

func (c httpClient) get(ctx context.Context, method string, rawURL string, out interface{}) error {
	rawURL = c.resolveURL(rawURL)
	req, err := http.NewRequestWithContext(ctx, method, rawURL, nil)
	if err != nil {
		return errors.New("creating http request failed").
			WithTag("method", method).
			WithTag("url", rawURL).
			Wrap(err)
	}
	return c.Do(req, out)
}

func (c httpClient) post(ctx context.Context, method string, rawURL string, out, in interface{}) error {
	rawURL = c.resolveURL(rawURL)
	body, err := c.encode(in)
	if err != nil {
		return errors.New("encoding http body failed").
			WithTag("method", method).
			WithTag("url", rawURL).
			Wrap(err)
	}

	req, err := http.NewRequestWithContext(ctx, method, rawURL, bytes.NewBuffer(body))
	if err != nil {
		return errors.New("creating http request failed").
			WithTag("method", method).
			WithTag("url", rawURL).
			Wrap(err)
	}
	req.Header.Set("Content-Type", "application/json")

	return c.Do(req, out)
}

func (c httpClient) resolveURL(v string) string {
	for _, r := range c.resolvers {
		v = r(v)
	}
	return v
}

func (c httpClient) Do(req *http.Request, out interface{}) error {
	for _, h := range c.headers {
		h(req.Header)
	}

	res, err := c.RoundTripper.RoundTrip(req)
	if err != nil {
		return errors.New("http request failed").
			WithTag("url", req.URL.String()).
			Wrap(err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return errors.New("http request failed").
			WithTag("url", req.URL.String()).
			WithTag("status", res.Status).
			WithTag("code", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.New("reading http response failed").
			WithTag("url", req.URL.String()).
			Wrap(err)
	}

	if out == nil {
		return nil
	}

	if err := c.decode(body, out); err != nil {
		return errors.New("decoding http response failed").
			WithTag("url", req.URL.String()).
			Wrap(err)
	}
	return nil
}
