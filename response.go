package fn

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
)

type response struct {
	Body            string            `json:"body"`
	Headers         map[string]string `json:"headers"`
	StatusCode      int               `json:"statusCode"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`

	httpHeaders http.Header   `json:"-"`
	bytes       *bytes.Buffer `json:"-"`
}

func newResponse() *response {
	return &response{
		Headers:     make(map[string]string),
		bytes:       new(bytes.Buffer),
		httpHeaders: make(http.Header),
	}
}
func (r *response) Header() http.Header {
	return r.httpHeaders
}

func (r *response) Write(b []byte) (int, error) {
	return r.bytes.Write(b)
}
func (r *response) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
}

func jsonResponse(code int, data interface{}) (*response, error) {
	res := &response{
		Headers: map[string]string{
			"content-type": "application/json",
		},
		StatusCode: code,
	}
	if data == nil {
		res.Body = "{}"
	} else if v, ok := data.(string); ok {
		res.Body = v
	} else {
		b, _ := json.Marshal(data)
		res.Body = string(b)
	}
	return res, nil
}

func stringResponse(data string) (*response, error) {
	res := &response{
		Headers: map[string]string{
			"content-type": "text/plain",
		},
		StatusCode: 200,
		Body:       data,
	}
	return res, nil
}

func (r *response) Wrap() {
	if r.httpHeaders != nil {
		if r.Headers == nil {
			r.Headers = make(map[string]string)
		}
		for k, vs := range r.httpHeaders {
			for _, v := range vs {
				r.Headers[k] = v
			}
		}
	}

	if r.bytes != nil && r.bytes.Len() > 0 {
		r.Body = base64.StdEncoding.EncodeToString(r.bytes.Bytes())
		r.IsBase64Encoded = true
	}
}
