package fastshot

import (
	"bryanchen463/cli_test/utils"
	"crypto/tls"
	"fmt"
	"net/http"

	fastshot "github.com/opus-domini/fast-shot"
)

var client = fastshot.NewClient(utils.Url).Build()

func init() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Create an HTTP client with the custom transport
	client := &http.Client{Transport: tr}
}

func POST(url string, payload map[string]interface{}) error {
	_, err := client.POST("").
		BodyJSON(payload).
		Send()

	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	return nil
	// Process response...
}

func GET(url string, payload map[string]interface{}) error {

	resp, err := client.GET("").
		BodyJSON(payload).
		Send()

	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	fmt.Println(resp.StatusCode())
	return nil
	// Process response...
}
