package heimdallcli

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gojek/heimdall/v7/hystrix"
	"gopkg.in/h2non/gentleman.v2/plugins/headers"
)

type myClient struct {
	client http.Client
}

func NewMyClient() *myClient {
	return &myClient{client: http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	},
	}
}

func (c *myClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

var client *hystrix.Client

func init() {
	myClient := NewMyClient()
	client = hystrix.NewClient(
		hystrix.WithHystrixTimeout(1100*time.Millisecond),
		hystrix.WithMaxConcurrentRequests(100),
		hystrix.WithErrorPercentThreshold(20),
		hystrix.WithHTTPClient(myClient),
	)
}

func Get(url string) error {

	header := http.Header{}
	headers.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Get(url, header)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if string(body) != "Hello, World 👋!" {
		return errors.New("response error")
	}
	return nil
}

func Post(url string, payload string) error {

	header := http.Header{}
	headers.Set("Content-Type", "application/x-www-form-urlencoded")
	postData := []byte(fmt.Sprintf(`name=%s`, payload))
	response, err := client.Post(url, bytes.NewBuffer(postData), header)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if string(body) != payload {
		return errors.New("response error " + payload + ": " + string(body))
	}
	return nil
}
