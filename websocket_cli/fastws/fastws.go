package fastwscli

// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/dgrr/fastws"
)

var conn *fastws.Conn
var buff []byte

func Init(addr string) error {
	var err error
	conn, err = fastws.DialTLS(addr, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return err
	}
	buff = make([]byte, 1024)
	return nil
}

func Clean() {
	conn.Close()
}

func SendReacv(message string) (int64, error) {
	conn.WriteMessage(fastws.ModeText, []byte(message))

	_, m, err := conn.ReadMessage(buff[:0])
	if err != nil {
		return 0, err
	}
	if string(m) != message {
		return 0, fmt.Errorf("received unexpected message, %s, %s", m, message)
	}
	return 0, nil
}

func Receive() (int64, error) {
	_, data, err := conn.ReadMessage(buff[:0])
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
