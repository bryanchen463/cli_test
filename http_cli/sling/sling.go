package slingcli

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/dghubble/sling"
)

var client *sling.Sling

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

func init() {
	client = sling.New()
	client.Client(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})
}

func Get(url string) error {
	client = sling.New()
	client.Client(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return err
	}
	resp, err := client.Do(req, nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if string(body) != "Hello, World 👋!" {
		return errors.New("response error")
	}
	return nil
}

func Post(url string, payload string) error {
	client = sling.New()
	client.Client(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})
	req := client.Post(url)
	req.Add("Content-Type", "application/x-www-form-urlencoded")
	postData := []byte(fmt.Sprintf(`name=%s`, payload))
	req.Body(bytes.NewBuffer(postData))
	r, err := req.Request()
	if err != nil {
		return err
	}
	resp, err := client.Do(r, nil, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if string(body) != payload {
		return errors.New("response error")
	}
	return nil
}
