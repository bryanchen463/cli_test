package cleanhttpcli

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/go-cleanhttp"
)

func Get(url string) error {
	cli := cleanhttp.DefaultPooledClient()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	_, err = cli.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func POST(url string, payload string) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
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
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("POST请求发送失败:", err)
		return err
	}
	defer resp.Body.Close()

	// 读取并处理响应
	if resp.Status == "200 OK" {
		// 在这里处理成功响应的逻辑
		_, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		// fmt.Println("body:", string(body))
	} else {
		// fmt.Println("POST请求失败:", resp.Status)
		// 在这里处理失败响应的逻辑
	}
	return nil
}
