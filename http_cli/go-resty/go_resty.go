package gorestycli

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

var client *resty.Client

func init() {
	client = resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
}

func Get(url string) error {
	resp, err := client.R().Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode())
	}
	if string(resp.Body()) != "Hello, World 👋!" {
		return errors.New("response error")
	}
	return nil
}

func Post(url string, payload string) error {
	req := client.R()
	req.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	postData := []byte(fmt.Sprintf(`name=%s`, payload))
	req.SetBody(postData)
	resp, err := req.Post(url)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}
	body := resp.Body()
	if string(body) != payload {
		return errors.New("response error")
	}
	return nil
}
