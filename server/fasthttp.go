package main

import (
	"github.com/dgrr/websocket"
	"github.com/valyala/fasthttp"
)

func fastHttpWs() {
	ws := websocket.Server{}
	ws.HandleData(OnMessage)

	fasthttp.ListenAndServeTLS(":8081", "../nginx/cert/server.crt", "../nginx/cert/server.key", ws.Upgrade)
}

func OnMessage(c *websocket.Conn, isBinary bool, data []byte) {
	c.Write(data)
}
