package requester

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	logger "github.com/sirupsen/logrus"
)

type RequestValidatorFunc func(*http.Request) error
type ClientOption func(*Client)

type Client struct {
	httpClient        *http.Client
	defaultOptions    []RequestOption
	requestValidators []RequestValidatorFunc
}

func (c *Client) Do(opts ...RequestOption) (*Response, error) {
	if opts == nil {
		opts = []RequestOption{}
	}
	opts = append(c.defaultOptions, opts...)
	req, err := NewRequest(opts...)
	if err != nil {
		return nil, err
	}

	for _, fn := range c.requestValidators {
		if err := fn(req.request); err != nil {
			return nil, err
		}
	}

	req.reqLogger = logger.StandardLogger().Writer()
	log.SetOutput(req.reqLogger)
	if req.reqLogger != nil {
		dump, err := httputil.DumpRequestOut(req.request, true)
		if err != nil {
			return nil, err
		}

		_, err = req.reqLogger.Write(dump)
		if err != nil {
			return nil, err
		}
	}

	httpResp, err := c.httpClient.Do(req.request)
	if err != nil {
		return nil, err
	}

	if req.respLogger != nil {
		dump, err := httputil.DumpResponse(httpResp, true)
		if err != nil {
			return nil, err
		}

		_, err = req.respLogger.Write(dump)
		if err != nil {
			return nil, err
		}
	}

	resp := &Response{
		Response: httpResp,
	}

	if err != nil {
		return resp, statusError{CodeUnknown, CodeUnknown, err.Error()}
	} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return resp, statusError{CodeBadResponseCode, resp.StatusCode, fmt.Sprintf("bad status code (%v)", resp.StatusCode)}
	}

	return resp, nil
}

func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		defaultOptions:    []RequestOption{},
		requestValidators: []RequestValidatorFunc{},
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

func WithDefaultOptions(opts ...RequestOption) ClientOption {
	return func(c *Client) { c.defaultOptions = append(c.defaultOptions, opts...) }
}

func WithRequestValidation(fn RequestValidatorFunc) ClientOption {
	return func(c *Client) { c.requestValidators = append(c.requestValidators, fn) }
}
