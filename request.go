package requester

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
)

type RequestOption func(*Request) error

type Request struct {
	*http.Request
	method  string
	host    string
	path    string
	url     string
	header  http.Header
	query   url.Values
	body    interface{}
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
		return nil, statusError{CodeMissingURL, "no host/url defined"}
	}

	var body io.Reader
	if r.body != nil {
		switch v := r.body.(type) {
		case io.Reader:
			body = v
		case string:
			body = bytes.NewBuffer([]byte(v))
		case []byte:
			body = bytes.NewBuffer(v)
		default:
			return nil, statusError{CodeInvalidBody, "invalid body type"}
		}
	}

	uri := r.url
	if uri == "" {
		uri = r.host + r.path
	}

	if req, err := http.NewRequest(r.method, uri, body); err != nil {
		return nil, statusError{CodeUnknown, err.Error()}
	} else {
		r.Request = req
	}

	for k, v := range r.header {
		r.Request.Header[k] = v
	}

	if len(r.query) > 0 {
		if len(r.Request.URL.RawQuery) > 0 {
			r.Request.URL.RawQuery += "&"
		}

		r.Request.URL.RawQuery += r.query.Encode()
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

func WithHeader(header http.Header) RequestOption {
	return func(r *Request) error {
		for k, v := range header {
			r.header[k] = v
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
			return statusError{CodeInvalidForm, "invalid form type, must be of type url.Values"}
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
