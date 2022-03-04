package requester

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
)

type RequestOption func(*Request) error

type Request struct {
	request    *http.Request
	ctx        context.Context
	method     string
	host       string
	path       string
	url        string
	header     http.Header
	query      url.Values
	body       interface{}
	reqLogger  io.Writer
	respLogger io.Writer
}

func NewRequest(opts ...RequestOption) (*Request, error) {
	r := &Request{
		header: http.Header{},
		query:  url.Values{},
	}

	if opts != nil {
		for _, opt := range opts {
			if err := opt(r); err != nil {
				return nil, err
			}
		}
	}

	if r.method == "" {
		r.method = http.MethodGet
	}
	if r.host == "" && r.url == "" {
		return nil, statusError{CodeMissingURL, CodeUnknown, "no host/url defined"}
	}

	var body io.Reader
	if r.body != nil {
		switch v := r.body.(type) {
		case io.Reader:
			body = v
		case *bytes.Buffer:
			body = v
		case string:
			body = bytes.NewBuffer([]byte(v))
		case []byte:
			body = bytes.NewBuffer(v)
		default:
			return nil, statusError{CodeInvalidBody, CodeUnknown, "invalid body type"}
		}
	}

	uri := r.url
	if uri == "" {
		uri = r.host + r.path
	}

	if r.ctx == nil {
		r.ctx = context.Background()
	}

	if req, err := http.NewRequestWithContext(r.ctx, r.method, uri, body); err != nil {
		return nil, statusError{CodeUnknown, CodeUnknown, err.Error()}
	} else {
		r.request = req
	}

	for k, v := range r.header {
		r.request.Header[k] = v
	}

	if len(r.query) > 0 {
		if len(r.request.URL.RawQuery) > 0 {
			r.request.URL.RawQuery += "&"
		}

		r.request.URL.RawQuery += r.query.Encode()
	}

	return r, nil
}

func WithMethod(method string) RequestOption {
	return func(r *Request) error {
		r.method = method
		return nil
	}
}

func WithHost(host string) RequestOption {
	return func(r *Request) error {
		r.host = host
		return nil
	}
}

func WithPath(path string) RequestOption {
	return func(r *Request) error {
		r.path = path
		return nil
	}
}

func WithURL(url string) RequestOption {
	return func(r *Request) error {
		r.url = url
		return nil
	}
}

func WithHeader(header map[string]string) RequestOption {
	return func(r *Request) error {
		for k, v := range header {
			r.header.Set(k, v)
		}
		return nil
	}
}

func WithQuery(query url.Values) RequestOption {
	return func(r *Request) error {
		for k, v := range query {
			r.query[k] = v
		}
		return nil
	}
}

func WithBody(v interface{}) RequestOption {
	return func(r *Request) error {
		r.body = v
		return nil
	}
}

func WithForm(v interface{}) RequestOption {
	return func(r *Request) error {
		q, ok := v.(url.Values)
		if !ok {
			return statusError{CodeInvalidForm, CodeUnknown, "invalid form type, must be of type url.Values"}
		}
		r.body = bytes.NewBuffer([]byte(q.Encode()))
		return nil
	}
}

func WithJSON(v interface{}) RequestOption {
	return func(r *Request) error {
		buf := new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(v); err != nil {
			return err
		}
		r.header["Content-Type"] = []string{"application/json"}
		r.body = buf
		return nil
	}
}

func WithXML(v interface{}) RequestOption {
	return func(r *Request) error {
		buf := new(bytes.Buffer)
		if err := xml.NewEncoder(buf).Encode(v); err != nil {
			return err
		}
		r.header["Content-Type"] = []string{"application/xml"}
		r.body = buf
		return nil
	}
}

func WithContext(ctx context.Context) RequestOption {
	return func(r *Request) error {
		r.ctx = ctx
		return nil
	}
}

func WithRequestLogger(w io.Writer) RequestOption {
	return func(r *Request) error {
		r.reqLogger = w
		return nil
	}
}

func WithResponseLogger(w io.Writer) RequestOption {
	return func(r *Request) error {
		r.respLogger = w
		return nil
	}
}
