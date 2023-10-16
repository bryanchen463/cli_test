package gentleman

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"

	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/body"
	gtls "gopkg.in/h2non/gentleman.v2/plugins/tls"
	"gopkg.in/h2non/gentleman.v2/plugins/transport"
)

var (
	postCli = gentleman.New().Post()
	getCli  = gentleman.New().Get()
)

func init() {
	getCli = gentleman.New().Post()
	getCli.Use(transport.Set(http.DefaultTransport))
	getCli.Use(gtls.Config(&tls.Config{
		InsecureSkipVerify: true,
	}))

	// Set a new header field
	getCli.SetHeader("Client", "gentleman")

	postCli.Use(transport.Set(http.DefaultTransport))
	postCli.Use(gtls.Config(&tls.Config{
		InsecureSkipVerify: true,
	}))

	// Set a new header field
	postCli.SetHeader("postClient", "gentleman")
	postCli.SetHeader("Content-Type", "application/x-www-form-urlencoded")

}

func Get(url string) error {
	// Create a new client
	// Create a custom TLS configuration to skip certificate verification
	// tlsConfig := &tls.Config{
	// 	InsecureSkipVerify: true,
	// }
	getCli = gentleman.New().Get()
	getCli.Use(transport.Set(http.DefaultTransport))
	getCli.Use(gtls.Config(&tls.Config{
		InsecureSkipVerify: true,
	}))

	// Set a new header field
	getCli.SetHeader("Client", "gentleman")
	getCli.URL(url)

	// Perform the request
	res, err := getCli.Send()
	if err != nil {
		fmt.Printf("Request error: %s\n", err)
		return err
	}
	if !res.Ok {
		fmt.Printf("Invalid server response: %d\n", res.StatusCode)
		return err
	}
	if "Hello, World 👋!" != string(res.Bytes()) {
		return errors.New("response error" + string(res.Bytes()))
	}
	return nil
}

func Post(url string, payload string) error {

	postCli = gentleman.New().Post()
	postCli.Use(gtls.Config(&tls.Config{
		InsecureSkipVerify: true,
	}))
	postCli.URL(url)
	postCli.Use(body.String(fmt.Sprintf("name=%s", payload)))

	// Perform the request
	res, err := postCli.Send()
	if err != nil {
		fmt.Printf("Request error: %s\n", err)
		return err
	}
	if !res.Ok {
		fmt.Printf("Invalid server response: %d\n", res.StatusCode)
		return errors.New("status code error")
	}
	if string(res.Bytes()) != payload {
		return errors.New("response error" + string(res.Bytes()))
	}
	return nil
}
