package gorillawebsocketclient

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var conn *websocket.Conn

func Init(addr string) error {
	dailer := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	var err error
	conn, _, err = dailer.Dial(addr, nil)
	if err != nil {
		return err
	}
	return nil
}

func Clean() {
	conn.Close()
}

func SendReacv(message string) (int64, error) {
	conn.SetWriteDeadline(time.Now().Add(time.Second))
	conn.SetReadDeadline(time.Now().Add(time.Second))
	now := time.Now()
	conn.WriteMessage(websocket.TextMessage, []byte(message))
	_, m, err := conn.ReadMessage()
	if err != nil {
		return 0, err
	}
	if string(m) != message {
		return 0, fmt.Errorf("received unexpected message, %s, %s", m, message)
	}
	return int64(time.Since(now).Microseconds()), nil
}

func Receive() (int64, error) {
	_, data, err := conn.ReadMessage()
	if err != nil {
		log.Println("receive:", err)
		return 0, err
	}
	// now := time.Now().UnixMicro()
	if string(data) == "finish" {
		return 0, nil
	}
	// var message common.Message
	// err = json.Unmarshal(data, &message)
	// if err != nil {
	// 	return 0, err
	// }

	return 0, nil
}
