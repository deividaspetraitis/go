package http

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/deividaspetraitis/go/errors"
	"github.com/deividaspetraitis/go/log"
)

// userAgent is the default user agent.
var userAgent = fmt.Sprintln("go")

// RequestOption is modifier of a request.
type RequestOption interface {
	apply(*http.Request)
}

// newRequestOption constructs a new requestOption.
func newRequestOption(fn func(r *http.Request)) *requestOption {
	return &requestOption{applyFn: fn}
}

// requestOption is an implementation of RequestOption.
type requestOption struct {
	applyFn func(r *http.Request)
}

// apply implements RequestOption.
func (o *requestOption) apply(r *http.Request) {
	o.applyFn(r)
}

// WithParameter constructs RequestOption to add arbitrary query parameters to the request.
func WithQueryParam(key, value string) RequestOption {
	return newRequestOption(func(r *http.Request) {
		q := r.URL.Query()
		q.Set(key, value)
		r.URL.RawQuery = q.Encode()
	})
}

// WithHeader constructs RequestOption to add arbitrary headers to the request.
func WithHeader(key, value string) RequestOption {
	return newRequestOption(func(r *http.Request) {
		r.Header.Add(key, value)
	})
}

// WithBearerToken constructs RequestOption to add Bearer authorisation token to the request.
func WithBearerToken(token string) RequestOption {
	return newRequestOption(func(r *http.Request) {
		r.Header.Add("Authorization", "Bearer "+token)
	})
}

// Client is a HTTP client capable to make requests to HTTP servers.
// It abstracts away lower details by providing high level interface.
type Client struct {
	url            *url.URL
	userAgent      string
	debug          bool
	http           *http.Client
	requestOptions []RequestOption
}

// NewClient constructs and returns new HTTP client instance.
func NewClient(domain string, opts ...RequestOption) (*Client, error) {
	u, err := url.Parse(domain)
	if err != nil {
		return nil, err
	}

	return &Client{
		url:            u,
		userAgent:      userAgent,
		debug:          false,
		http:           http.DefaultClient,
		requestOptions: opts,
	}, nil
}

// SetRequestOption replaces RequestOption's for the client reused across all client issued requests.
func (c *Client) SetRequestOption(r ...RequestOption) {
	c.requestOptions = r
}

// AddRequestOption adds RequestOption's to the client reused across all client issued requests.
func (c *Client) AddRequestOption(r ...RequestOption) {
	c.requestOptions = append(c.requestOptions, r...)
}

// RequestOptions returs RequestOption's used by the client.
func (c *Client) RequestOptions() []RequestOption {
	return c.requestOptions
}

// URI returns the absolute URL of the API with any path segments
// appended to the end.
func (c *Client) URI(path ...string) string {
	return (&url.URL{
		Scheme: c.url.Scheme,
		Host:   c.url.Host,
		Path:   c.url.Path + "/" + strings.Join(path, "/"),
	}).String()
}

// NewRequest returns a new HTTP request. If the payload is not nil it will be encoded as JSON.
func (c *Client) request(ctx context.Context, method, uri string, payload []byte, options ...RequestOption) (r *http.Request, err error) {
	r, err = http.NewRequestWithContext(ctx, method, uri, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	// Apply global client options
	for _, opt := range c.requestOptions {
		if opt != nil {
			opt.apply(r)
		}
	}

	// Apply request specific options
	for _, opt := range options {
		if opt != nil {
			opt.apply(r)
		}
	}

	return
}

// Request combines request and do, while also handling decoding of response
// payload.
func (c *Client) Request(ctx context.Context, method, uri string, v []byte, options ...RequestOption) (*http.Response, error) {
	uri = c.URI(uri)

	req, err := c.request(ctx, method, uri, v, options...)
	if err != nil {
		return nil, errors.Wrapf(err, "building request")
	}

	res, err := c.do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "sending request to %s", uri)
	}

	if c.debug {
		log.Printf("request to %s resulted in HTTP response code %d", req.URL.String(), res.StatusCode)
	}

	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return res, errors.Newf("request resulted in %d response code", res.StatusCode)
	}

	return res, nil
}

// do sends an HTTP request and returns an HTTP response, handling any context
// cancellations or timeouts.
func (c *Client) do(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	res, err := c.http.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			return nil, err
		}
	}

	return res, nil
}
