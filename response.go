package httpcli

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type response struct {
	resp *http.Response
	body []byte
}

func newResponse(resp *http.Response) (r *response) {
	r = &response{
		resp: resp,
	}
	return
}

var ErrNilResponse = errors.New("nil response")

func (resp *response) Response() *http.Response {
	if resp == nil {
		return nil
	}
	return resp.resp
}

func (resp *response) receive() (err error) {
	if resp == nil {
		err = ErrNilResponse
		return
	}
	if resp.body != nil {
		return
	}
	defer resp.resp.Body.Close()
	resp.body, err = ioutil.ReadAll(resp.resp.Body)
	return
}

func (resp *response) ReceiveBytes() (body []byte, err error) {
	if resp == nil {
		err = ErrNilResponse
		return
	}
	if resp.body == nil {
		if err = resp.receive(); err != nil {
			return
		}
	}
	body = resp.body
	return
}

func (resp *response) Bytes() (body []byte) {
	body, _ = resp.ReceiveBytes()
	return
}

func (resp *response) ReceiveString() (s string, err error) {
	if resp == nil {
		err = ErrNilResponse
		return
	}
	if resp.body == nil {
		if err = resp.receive(); err != nil {
			return
		}
	}
	s = string(resp.body)
	return
}

func (resp *response) String() (s string) {
	s, _ = resp.ReceiveString()
	return
}

func (resp *response) Json(v interface{}) (err error) {
	if resp == nil {
		err = ErrNilResponse
		return
	}
	if resp.body == nil {
		if err = resp.receive(); err != nil {
			return
		}
	}
	err = json.Unmarshal(resp.body, v)
	return
}

func (resp *response) Xml(v interface{}) (err error) {
	if resp == nil {
		err = ErrNilResponse
		return
	}
	if resp.body == nil {
		if err = resp.receive(); err != nil {
			return
		}
	}
	err = xml.Unmarshal(resp.body, v)
	return
}
func (resp *response) Status() string {
	return resp.resp.Status
}
func (resp *response) StatusCode() int {
	return resp.resp.StatusCode
}

// info return resposne detail with string
func (resp *response) Info() string {
	if resp.resp == nil {
		return ""
	}
	var (
		out      bytes.Buffer
		str      string
		response = resp.resp
	)
	if str = resp.String(); str == "" {
		return str
	}
	out.WriteString(fmt.Sprint(response.Proto, " ", resp.Status))
	if len(response.Header) > 0 {
		for name, values := range response.Header {
			for _, value := range values {
				out.WriteString(fmt.Sprintf("\n%s:%s", name, value))
			}
		}
	}
	//body
	out.WriteString(fmt.Sprint("\n\n", str))
	return out.String()
}
