package reqcli

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/imroc/req/v3"
)

var client *req.Client

func init() {
	client = req.C()
	client = client.EnableInsecureSkipVerify()
}

func Get(url string) error {
	req := client.Get(url)
	resp := req.Do()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
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
	req := client.Post(url)
	req.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	postData := []byte(fmt.Sprintf(`name=%s`, payload))
	req.SetBody(postData)
	resp := req.Do()
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
