package nhooyrcli

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"

	"nhooyr.io/websocket"
)

func sendReacv(conn *websocket.Conn, message string) (int, error) {
	ctx := context.Background()
	conn.Write(ctx, websocket.MessageText, []byte(message))
	_, m, err := conn.Read(ctx)
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

	c, _, err := websocket.Dial(context.Background(), addr, &websocket.DialOptions{
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
	defer c.Close(websocket.StatusNormalClosure, "finished")

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
