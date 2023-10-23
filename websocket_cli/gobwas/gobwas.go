package gobwascli

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

var conn net.Conn

func Init(addr string) error {
	dailer := ws.Dialer{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	var err error
	conn, _, _, err = dailer.Dial(context.Background(), addr)
	if err != nil {
		return err
	}
	return nil
}

func Clean() {
	conn.Close()
}

func SendReacv(message string) (int64, error) {
	// conn.SetWriteDeadline(time.Now().Add(time.Second))
	// conn.SetReadDeadline(time.Now().Add(time.Second))
	// now := time.Now()
	err := wsutil.WriteClientMessage(conn, ws.OpText, []byte(message))
	if err != nil {
		return 0, err
	}
	data, _, err := wsutil.ReadServerData(conn)
	if err != nil {
		return 0, err
	}
	if string(data) != message {
		return 0, fmt.Errorf("received unexpected message, %s, %s", string(data), message)
	}
	return 0, nil
}

func Receive() (int64, error) {
	data, _, err := wsutil.ReadServerData(conn)
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
