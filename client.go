package requester

import (
	"fmt"
	"net/http"
)

type ClientOption func(*Client)

type Client struct {
	httpClient     *http.Client
	defaultOptions []RequestOption
}

func (c *Client) Do(opts ...RequestOption) (*Response, error) {
	req, err := NewRequest(append(c.defaultOptions, opts...)...)
	if err != nil {
		return nil, err
	}

	httpResp, err := c.httpClient.Do(req.request)
	resp := &Response{
		Response: httpResp,
	}
	if err != nil {
		return nil, statusError{CodeUnknown, err.Error()}
	} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, statusError{CodeBadResponseStatus, fmt.Sprintf("bad status code (%v)", resp.StatusCode)}
	}

	return resp, nil
}

func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		defaultOptions: []RequestOption{},
	}

	if opts != nil {
		for _, opt := range opts {
			opt(c)
		}
	}

	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}

	return c
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) { c.httpClient = httpClient }
}

func WithDefaultHeader(header http.Header) ClientOption {
	return func(c *Client) { c.defaultOptions = append(c.defaultOptions, WithHeader(header)) }
}

func WithDefaultHost(host string) ClientOption {
	return func(c *Client) { c.defaultOptions = append(c.defaultOptions, WithHost(host)) }
}
