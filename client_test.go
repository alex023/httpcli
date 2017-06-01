package httpcli

import (
	"testing"
)

func TestGet(t *testing.T) {
	req := Get("http://www.baidu.com")
	_, err := req.Response()
	if err != nil || len(req.Info()) == 0 {
		t.Error(err)
	}
}
func TestRequest_InsecureTLS(t *testing.T) {
	req := Get("https://www.google.com").WithTLS()
	_, err := req.Response()
	if err != nil || len(req.Info()) == 0 {
		t.Error(err)
	}
}
func TestPost(t *testing.T) {
	req := Post("https://baidu.com").WithTLS().
		WithParam("wd", "查询")
	//t.Log(req.Info())
	resp, err := req.Response()
	if err != nil || len(resp.Info()) == 0 {
		t.Error(err)
	}
}
