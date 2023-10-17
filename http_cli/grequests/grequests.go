package grequestscli

import (
	"errors"
	"log"

	"github.com/levigross/grequests"
)

var session *grequests.Session

func init() {
	ro := &grequests.RequestOptions{InsecureSkipVerify: true}
	session = grequests.NewSession(ro)
}

func Get(url string) error {
	resp, err := session.Get(url, nil)
	// You can modify the request by passing an optional RequestOptions struct

	if err != nil {
		log.Fatalln("Unable to make request: ", err)
	}
	defer resp.Close()

	if resp.String() != "Hello, World 👋!" {
		return errors.New("response error")
	}
	return nil
}

func Post(url string, payload string) error {
	ro := &grequests.RequestOptions{
		Data: map[string]string{
			"name": payload,
		},
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
	}
	resp, err := session.Post(url, ro)
	// You can modify the request by passing an optional RequestOptions struct

	if err != nil {
		log.Fatalln("Unable to make request: ", err)
	}
	defer resp.Close()

	if resp.String() != payload {
		return errors.New("response error " + payload)
	}
	return nil
}
