package nbiocli

/*
import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/lesismal/llib/std/crypto/tls"
	"github.com/lesismal/nbio/nbhttp"
	"github.com/lesismal/nbio/nbhttp/websocket"
)

var (
	connected    uint64 = 0
	success      uint64 = 0
	failed       uint64 = 0
	totalSuccess uint64 = 0
	totalFailed  uint64 = 0
	conn         *websocket.Conn
)

func newUpgrader() *websocket.Upgrader {
	u := websocket.NewUpgrader()

	u.OnClose(func(c *websocket.Conn, err error) {
		fmt.Println("OnClose:", c.RemoteAddr().String(), err)
	})

	return u
}

func Init(addr string) error {

	engine := nbhttp.NewEngine(nbhttp.Config{})
	err := engine.Start()
	if err != nil {
		fmt.Printf("nbio.Start failed: %v\n", err)
		return err
	}

	u := url.URL{Scheme: "wss", Host: addr, Path: "/wss"}
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	dialer := &websocket.Dialer{
		Engine:          engine,
		Upgrader:        newUpgrader(),
		DialTimeout:     time.Second * 3,
		TLSClientConfig: tlsConfig,
	}

	conn, _, err = dialer.Dial(u.String(), nil)
	return err
}

func Clean() {
	conn.CloseAndClean()
}

func SendRead(message string) (int64, error) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		return 0, err
	}
	_, err = conn.Read(buff)
	if err != nil {
		return 0, err
	}
	if string(buff) != message {
		return 0, fmt.Errorf("received unexpected message, %s, %s", string(buff), message)
	}
	return 0, nil
}

func Receive() (int64, error) {
	_, err := tlsConn.Read(buff)
	if err != nil {
		return 0, err
	}
	if err != nil {
		log.Println("receive:", err)
		return 0, err
	}
	if string(buff) == "finish" {
		return 0, nil
	}
	return 0, nil
}
*/
