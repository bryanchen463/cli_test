package goretryablenhttpcli

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

var (
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Create an HTTP client with the custom transport
	client = &http.Client{Transport: tr}
)

func Get(url string) error {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10
	retryClient.HTTPClient = client

	standardClient := retryClient.StandardClient()
	resp, err := standardClient.Get(url)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if string(body) != "Hello, World 👋!" {
		return errors.New("response error")
	}
	return nil
}

func POST(url string, payload string) error {

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10
	retryClient.HTTPClient = client

	standardClient := retryClient.StandardClient()

	// Create an HTTP client with the custom transport
	postData := []byte(fmt.Sprintf(`name=%s`, payload))
	// 创建一个请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postData))
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 发送请求
	resp, err := standardClient.Do(req)
	if err != nil {
		fmt.Println("POST请求发送失败:", err)
		return err
	}
	defer resp.Body.Close()

	// 读取并处理响应
	if resp.Status == "200 OK" {
		// 在这里处理成功响应的逻辑
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if string(body) != payload {
			return errors.New("response error")
		}
		// fmt.Println("body:", string(body))
	} else {
		// fmt.Println("POST请求失败:", resp.Status)
		// 在这里处理失败响应的逻辑
	}
	return nil
}
