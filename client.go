//Package httpcli is toolkit for simplification of http request
package httpcli

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Proto struct {
}

const (
	HTTP1 = "HTTP/1.1"
	HTTP2 = "HTTP/2.0"
)

// Params represents the request params.
type Params map[string]string

// Client is a toolkit base on http.Client
type Client struct {
	url       string
	urlEncode bool
	params    Params
	req       *http.Request
	resp      *response
	body      []byte
	client    http.Client
}

// Request return the raw *http.Client.
func (client *Client) Request() *http.Request {
	return client.req
}

// WithTLS insecure the https.
func (client *Client) WithTLS() *Client {
	client.client.Transport = &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	return client
}

// WithParam set single param to the request.
func (client *Client) WithParam(k, v string) *Client {
	client.params[k] = v
	return client
}

// WithParams set multiple params to the request.
func (client *Client) WithParams(params Params) *Client {
	for k, v := range params {
		client.params[k] = v
	}
	return client
}

// WithHeader set the request header.
func (client *Client) WithHeader(k, v string) *Client {
	client.req.Header.Set(k, v)
	return client
}

// WithHeaders set multiple headers.
func (client *Client) WithHeaders(params Params) *Client {
	for k, v := range params {
		client.req.Header.Set(k, v)
	}
	return client
}

// WithBody set  request body,support string and []byte.
func (client *Client) WithBody(body interface{}) *Client {
	switch v := body.(type) {
	case []byte:
		bf := bytes.NewBuffer(v)
		client.req.Body = ioutil.NopCloser(bf)
		client.req.ContentLength = int64(len(v))
		client.body = v
	case string:
		bf := bytes.NewBufferString(v)
		client.req.Body = ioutil.NopCloser(bf)
		client.req.ContentLength = int64(len(v))
		client.body = []byte(v)
	}
	return client
}

// GetBody return the request body.
func (client *Client) GetBody() []byte {
	if client.body == nil && client.req.Method == "POST" {
		return []byte(client.getParamBody())
	}
	return client.body
}

func (client *Client) getParamBody() string {
	if len(client.params) == 0 {
		return ""
	}
	var buf bytes.Buffer
	for k, v := range client.params {
		if client.urlEncode {
			k = url.QueryEscape(k)
			v = url.QueryEscape(v)
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(v)
		buf.WriteByte('&')
	}
	p := buf.String()
	p = p[0 : len(p)-1]
	return p
}

func (client *Client) buildGetUrl() string {
	ret := client.url
	if p := client.getParamBody(); p != "" {
		if strings.Index(client.url, "?") != -1 {
			ret += "&" + p
		} else {
			ret += "?" + p
		}
	}
	return ret
}

func (client *Client) WithJson(content string) *Client {
	client.WithBody(content)
	client.WithHeader("Content-Type", "application/json")
	client.urlEncode = false
	return client
}

func (client *Client) setParamBody() {
	if client.urlEncode {
		client.WithHeader("Content-Type", "application/x-www-form-urlencoded")
	}
	client.WithBody(client.getParamBody())
}

// GetUrl return the url of the request.
func (client *Client) GetUrl() string {
	if client.req.Method == "GET" {
		return client.buildGetUrl() //GET method and did not send request yet.
	}
	return client.url
}

// Url set the request's url.
func (client *Client) Url(urlStr string) *Client {
	client.url = urlStr
	return client
}

// Response execute the request and get the response, return error if error happens.
func (client *Client) Response() (resp *response, err error) {
	if client.resp != nil {
		resp = client.resp
		return
	}
	err = client.Do()
	if err != nil {
		return
	}
	resp = client.resp
	return
}

// Undo let the request could be executed again.
func (client *Client) Undo() *Client {
	client.resp = nil
	return client
}

// Do just execute the request. return error if error happens.
func (client *Client) Do() (err error) {
	// handle request params
	destUrl := client.url
	if len(client.params) > 0 {
		switch client.req.Method {
		case "GET":
			destUrl = client.buildGetUrl()
		case "POST":
			client.setParamBody()
		}
	}
	// set url
	u, err := url.Parse(destUrl)
	if err != nil {
		return
	}
	client.req.URL = u
	respRaw, err := client.client.Do(client.req)
	if err != nil {
		return
	}
	resp := newResponse(respRaw)
	err = resp.receive()
	if err != nil {
		return
	}
	client.resp = resp
	return
}

// Get returns *Client with GET method.
func Get(url string) *Client {
	return newRequest(url, "GET", HTTP1)
}

// Post returns *Client with POST method.
func Post(url string) *Client {
	return newRequest(url, "POST", HTTP1)
}

// New return a Client with the underlying *http.Client.
func New(req *http.Request) *Client {
	return &Client{
		urlEncode: true,
		params:    Params{},
		req:       req,
	}
}

func newRequest(url, method, proto string) *Client {
	return &Client{
		url:       url,
		urlEncode: true,
		params:    Params{},
		req: &http.Request{
			Method:     method,
			Header:     make(http.Header),
			Proto:      proto,
			ProtoMajor: 1,
			ProtoMinor: 1,
		},
	}
}

//Info return the client request information both response if it is not nil.
func (client *Client) Info() string {
	var out bytes.Buffer
	out.WriteString(fmt.Sprint(client.req.Method, " ", client.GetUrl(), " ", client.req.Proto))

	for name, values := range client.req.Header {
		for _, value := range values {
			out.WriteString(fmt.Sprint("\n\n", name, ":", value))
		}
	}
	if len(client.body) > 0 {
		out.WriteString(fmt.Sprint("\n\n", string(client.body)))
	}
	if client.resp != nil {
		out.WriteString(fmt.Sprint("\n\n", client.resp.Info()))
	}
	return out.String()
}
