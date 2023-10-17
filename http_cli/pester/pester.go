package pestercli

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/sethgrid/pester"
)

var client *pester.Client

func init() {
	client = pester.New()
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
}

func Get(url string) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
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
