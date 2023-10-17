package gorillawebsocketclient

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func sendReacv(conn *websocket.Conn, message string) (int, error) {
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

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

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

	done := make(chan struct{})

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	seq := 0
	for seq < len(message) {
		select {
		case <-done:
			return nil
		case <-ticker.C:
			_, err := sendReacv(c, message[seq])
			if err != nil {
				log.Println("write:", err)
				return err
			}
			seq++
		}
	}
	return nil
}
