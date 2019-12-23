package requester

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
)

type decodeFunc func(io.Reader, interface{}) error

func jsonDecode(r io.Reader, v interface{}) error { return json.NewDecoder(r).Decode(v) }
func xmlDecode(r io.Reader, v interface{}) error  { return xml.NewDecoder(r).Decode(v) }

type Response struct {
	*http.Response
}

func (r *Response) decode(v interface{}, decode decodeFunc) error {
	defer r.Body.Close()

	return decode(r.Body, v)
}

func (r *Response) JSON(v interface{}) error {
	return r.decode(v, jsonDecode)
}

func (r *Response) XML(v interface{}) error {
	return r.decode(v, xmlDecode)
}
