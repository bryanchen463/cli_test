package nhooyrcli

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"nhooyr.io/websocket"
)

var conn *websocket.Conn

func Init(addr string) error {
	var err error
	conn, _, err = websocket.Dial(context.Background(), addr, &websocket.DialOptions{
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func Clean() {
	defer conn.Close(websocket.StatusNormalClosure, "finished")
}

func SendReacv(message string) (int64, error) {
	ctx := context.Background()
	// ctx, cancel := context.WithTimeout(ctx, time.Second)
	// defer cancel()
	// now := time.Now()
	conn.Write(ctx, websocket.MessageText, []byte(message))
	_, m, err := conn.Read(ctx)
	if err != nil {
		return 0, err
	}
	if string(m) != message {
		return 0, fmt.Errorf("received unexpected message, %s, %s", m, message)
	}
	// return int64(time.Since(now).Microseconds()), nil
	return 0, nil
}

func Receive() (int64, error) {

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	_, data, err := conn.Read(ctx)
	if err != nil {
		log.Println("receive:", err)
		return 0, err
	}
	if string(data) == "finish" {
		return 0, nil
	}
	// now := time.Now().UnixMicro()
	// if string(data) == "finish" {
	// 	return 0, nil
	// }
	// var message common.Message
	// err = json.Unmarshal(data, &message)
	// if err != nil {
	// 	return 0, err
	// }

	return 0, nil
}
