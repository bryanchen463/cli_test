package gorillawebsocketclient

import (
	"flag"
	"fmt"
	"log"
	"net/url"
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

func Start(addr string, message []string) {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: addr}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	seq := 0
	for seq < len(message) {
		select {
		case <-done:
			return
		case <-ticker.C:
			_, err := sendReacv(c, message[seq])
			if err != nil {
				log.Println("write:", err)
				return
			}
			seq++
		}
	}
}
