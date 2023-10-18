package gorillawebsocketclient

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func sendReacv(conn *websocket.Conn, message string) (int, error) {
	conn.SetWriteDeadline(time.Now().Add(time.Second))
	conn.SetReadDeadline(time.Now().Add(time.Second))
	conn.WriteMessage(websocket.TextMessage, []byte(message))
	_, m, err := conn.ReadMessage()
	if err != nil {
		return 0, err
	}
	if string(m) != message {
		return 0, fmt.Errorf("received unexpected message, %s, %s", m, message)
	}
	return len(m), nil
}

func Start(addr string, message []string) error {
	flag.Parse()
	log.SetFlags(0)

	dailer := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	c, _, err := dailer.Dial(addr, nil)
	if err != nil {
		return err
	}
	defer c.Close()

	seq := 0
	for _, m := range message {
		_, err := sendReacv(c, m)
		if err != nil {
			log.Println("write:", err)
			return err
		}
		seq++
	}
	return nil
}
