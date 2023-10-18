package fastwscli

// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/dgrr/fastws"
)

func sendReacv(conn *fastws.Conn, message string) (int, error) {
	conn.WriteMessage(fastws.ModeText, []byte(message))
	buff := make([]byte, 0, len(message))
	_, m, err := conn.ReadMessage(buff)
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
	c, err := fastws.DialTLS(addr, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return err
	}
	defer c.Close()
	c.ReadTimeout = time.Second
	c.WriteTimeout = time.Second
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for _, m := range message {

		_, err := sendReacv(c, m)
		if err != nil {
			log.Println("write:", err)
			return err
		}
	}
	return nil
}
