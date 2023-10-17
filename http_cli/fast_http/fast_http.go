package fasthttpcli

import (
	"crypto/tls"
	"errors"
	"fmt"

	"github.com/valyala/fasthttp"
)

var (
	getCli  = fasthttp.Client{}
	postCli = fasthttp.Client{}
)

func init() {
	getCli.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	postCli.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
}

func Get(url string) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod("GET")
	req.SetRequestURI(url)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := getCli.Do(req, resp)
	if resp.StatusCode() != 200 {
		return errors.New("status code is not equal 200")
	}
	if "Hello, World 👋!" != string(resp.Body()) {
		return errors.New("response error" + string(resp.Body()))
	}
	return err
}

func Post(url string, payload string) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/x-www-form-urlencoded")
	req.SetBody([]byte(fmt.Sprintf("name=%s", payload)))
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := postCli.Do(req, resp)
	// fmt.Println(string(resp.Body()))
	if string(resp.Body()) != payload {
		return errors.New("response error")
	}
	return err
}
