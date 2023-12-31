package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/antlabs/quickws"
)

type echoHandler struct{}

func (e *echoHandler) OnOpen(c *quickws.Conn) {
	// fmt.Println("OnOpen:", c.NetConn().RemoteAddr())
}

func (e *echoHandler) OnMessage(c *quickws.Conn, op quickws.Opcode, msg []byte) {
	// fmt.Printf("OnMessage: %s, %v\n", msg, op)
	if err := c.WriteTimeout(op, msg, 3*time.Second); err != nil {
		fmt.Println("write fail:", err)
	}
}

func (e *echoHandler) OnClose(c *quickws.Conn, err error) {
	fmt.Printf("OnClose: %v\n", err)
}

// echo测试服务
func echo(w http.ResponseWriter, r *http.Request) {
	c, err := quickws.Upgrade(w, r, quickws.WithServerReplyPing(),
		// quickws.WithServerDecompression(),
		// quickws.WithServerIgnorePong(),
		quickws.WithServerCallback(&echoHandler{}),
		quickws.WithServerReadTimeout(5*time.Second),
	)
	if err != nil {
		fmt.Println("Upgrade fail:", err)
		return
	}

	c.StartReadLoop()
}

func quickwsServer() {
	http.HandleFunc("/ws/echo", echo)

	http.ListenAndServeTLS(":8080", "../nginx/cert/server.crt", "../nginx/cert/server.key", nil)
}
