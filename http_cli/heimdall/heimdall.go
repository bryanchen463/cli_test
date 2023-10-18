package heimdallcli

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gojek/heimdall/v7/hystrix"
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
		hystrix.WithHTTPClient(myClient),
	)
}

func Get(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return err
	}
	response, err := client.Do(req)
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

	postData := []byte(fmt.Sprintf(`name=%s`, payload))
	// 设置请求头
	// Create an HTTP client with the custom transport
	// 创建一个请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postData))
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("status: %d", response.StatusCode)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if string(body) != payload {
		return errors.New("response error " + payload + ": " + string(body))
	}
	return nil
}
