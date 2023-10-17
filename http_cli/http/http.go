package http

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Create an HTTP client with the custom transport
	client = &http.Client{Transport: tr}
)

func POST(url string, payload string) error {

	postData := []byte(fmt.Sprintf(`name=%s`, payload))

	// Create an HTTP client with the custom transport
	// 创建一个请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postData))
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("POST请求发送失败:", err)
		return err
	}
	defer resp.Body.Close()

	// 读取并处理响应
	if resp.Status == "200 OK" {
		// 在这里处理成功响应的逻辑
		body, err := io.ReadAll(resp.Body)
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

func GET(url string) error {

	// URL with an unsafe certificate (replace with your own URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return err
	}
	// Make a GET request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	defer resp.Body.Close()

	// Read and print the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return err
	}
	if "Hello, World 👋!" != string(body) {
		return errors.New("response error")
	}
	return nil
}
