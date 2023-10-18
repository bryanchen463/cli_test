package gobwascli

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/gobwas/ws"
)

func sendReacv(conn net.Conn, message string) (int, error) {
	conn.Write([]byte(message))
	buff := make([]byte, 0, len(message))
	start := 0
	for {
		n, err := conn.Read(buff[start:])
		if err != nil {
			return 0, err
		}
		if start+n == len(buff)-1 {
			if string(buff) != message {
				return 0, fmt.Errorf("received unexpected message, %s, %s", string(buff), message)
			}
			return 0, errors.New("received unexpected message")
		}
	}
}

func Start(addr string, message []string) error {
	flag.Parse()
	log.SetFlags(0)

	dailer := ws.Dialer{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	c, _, _, err := dailer.Dial(context.Background(), addr)
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
