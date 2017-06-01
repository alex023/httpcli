# httpcli
httpcli is better choice for http request.

[![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/alex023/httpcli)](https://goreportcard.com/report/github.com/alex023/httpcli)
[![GoDoc](https://godoc.org/github.com/alex023/httpcli?status.svg)](https://godoc.org/github.com/alex023/httpcli)
[![Build Status](https://travis-ci.org/alex023/httpcli.svg?branch=dev)](https://travis-ci.org/alex023/httpcli?branch=dev)
[![Coverage Status](https://coveralls.io/repos/github/alex023/httpcli/badge.svg?branch=dev)](https://coveralls.io/github/alex023/httpcli?branch=dev)

## Brief
implements a friendly API for net/http standard library.
# Features
- simpler Get、Post request
- decorate function param settings 
## Example
### Get a page

```golang
	req := httpcli.Get("https://www.google.com").WithTLS()
	response, err := req.Response()
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Println("result:",response.String())
	}
```
### Post a req
```golang
	req := Post("https://baidu.com").WithTLS().
    		WithParam("wd", "查询")
	response, err := req.Response()
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Println("result:",response.String())
	}
```