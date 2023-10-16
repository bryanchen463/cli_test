package gohttpclient

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"

	gohttpclient "github.com/bozd4g/go-http-client"
)

var (
	TlsClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
)

func Get(url string) error {
	ctx := context.Background()
	client := gohttpclient.New(url)

	response, err := client.Get(ctx, "/", gohttpclient.Option(gohttpclient.WithCustomHttpClient(TlsClient)))
	if err != nil {
		log.Fatalf("get error: %v", err)
	}
	if string(response.Body()) != "Hello, World 👋!" {
		return errors.New("response error")
	}
	return nil
}

func Post(url string, payload string) error {
	ctx := context.Background()
	client := gohttpclient.New(url)
	response, err := client.Post(ctx, "/", gohttpclient.Option(gohttpclient.WithCustomHttpClient(TlsClient)), gohttpclient.WithBody([]byte(fmt.Sprintf("name=%s", payload))), gohttpclient.WithHeader("Content-Type", "application/x-www-form-urlencoded"))
	if err != nil {
		log.Fatalf("post error: %v", err)
	}
	if string(response.Body()) != payload {
		return errors.New("response error")
	}
	return nil
}
