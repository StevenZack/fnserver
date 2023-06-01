package fn

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type (
	request struct {
		Version               string            `json:"version"`
		RouteKey              string            `json:"routeKey"`
		RawPath               string            `json:"rawPath"`
		RawQueryString        string            `json:"rawQueryString"`
		Cookies               []string          `json:"cookies"`
		PathParameters        map[string]string `json:"pathParameters"`
		IsBase64Encoded       bool              `json:"isBase64Encoded"`
		Body                  string            `json:"body"`
		RequestContext        requestContext    `json:"requestContext"`
		Headers               map[string]string `json:"headers"`
		QueryStringParameters map[string]string `json:"queryStringParameters"`
	}
	requestContext struct {
		AccountId    string      `json:"accountId"`
		ApiId        string      `json:"apiId"`
		DomainName   string      `json:"domainName"`
		DomainPrefix string      `json:"domainPrefix"`
		Http         httpContext `json:"http"`
		RequestId    string      `json:"requestId"`
		RouteKey     string      `json:"routeKey"`
		Stage        string      `json:"stage"`
		Time         string      `json:"time"`
		TimeEpoch    int64       `json:"timeEpoch"`
	}
	httpContext struct {
		Method    string `json:"method"`
		Path      string `json:"path"`
		Protocal  string `json:"protocal"`
		SourceIp  string `json:"sourceIp"`
		UserAgent string `json:"userAgent"`
	}
)

func (r *request) GetHeader(key string) string {
	if v, ok := r.Headers[key]; ok {
		return v
	}
	return r.Headers[strings.ToLower(key)]
}

func (r *request) GetMethod() string {
	return r.RequestContext.Http.Method
}

func (r *request) Language() string {
	if v, ok := r.Headers["language"]; ok {
		return v
	}
	if v, ok := r.Headers["Language"]; ok {
		return v
	}
	return r.Headers["accept-language"]
}

func (r *request) RawBody() ([]byte, error) {
	var data []byte
	var e error
	if r.IsBase64Encoded {
		data, e = base64.StdEncoding.DecodeString(r.Body)
		if e != nil {
			return nil, e
		}
	} else {
		data = []byte(r.Body)
	}

	return data, nil
}

func (e *request) String() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func (r *request) ToHttpRequest() (*http.Request, error) {
	body, e := r.RawBody()
	if e != nil {
		log.Println(e)
		return nil, e
	}

	url := "http://" + r.RequestContext.DomainName + r.RawPath
	if r.RawQueryString != "" {
		url += "?" + r.RawQueryString
	}
	req, e := http.NewRequest(r.GetMethod(), url, bytes.NewReader(body))
	if e != nil {
		log.Println(e)
		return nil, e
	}
	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}

	return req, nil
}
