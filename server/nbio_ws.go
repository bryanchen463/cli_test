package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/lesismal/llib/std/crypto/tls"
	"github.com/lesismal/nbio/nbhttp"
	"github.com/lesismal/nbio/nbhttp/websocket"
)

var (
	svr   *nbhttp.Server
	print = flag.Bool("print", false, "stdout output of echoed data")
)

func newUpgrader() *websocket.Upgrader {
	u := websocket.NewUpgrader()
	u.OnMessage(func(c *websocket.Conn, messageType websocket.MessageType, data []byte) {
		// echo
		c.WriteMessage(messageType, data)

		c.SetReadDeadline(time.Now().Add(nbhttp.DefaultKeepaliveTime))
	})
	u.OnClose(func(c *websocket.Conn, err error) {
		if *print {
			fmt.Println("OnClose:", c.RemoteAddr().String(), err)
		}
	})

	return u
}

func onWebsocket(w http.ResponseWriter, r *http.Request) {
	// time.Sleep(time.Second * 5)
	upgrader := newUpgrader()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	if *print {
		fmt.Println("OnOpen:", conn.RemoteAddr().String())
	}
}

func nbioWsServer() {
	flag.Parse()
	certFile, err := os.ReadFile("../nginx/cert/server.crt")
	if err != nil {
		log.Println(err)
		return
	}
	keyFile, err := os.ReadFile("../nginx/cert/server.key")
	if err != nil {
		log.Println(err)
		return
	}
	cert, err := tls.X509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("tls.X509KeyPair failed: %v", err)
	}
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	mux := &http.ServeMux{}
	mux.HandleFunc("/ws/echo", onWebsocket)

	svr = nbhttp.NewServer(nbhttp.Config{
		Network:   "tcp",
		AddrsTLS:  []string{"localhost:8888"},
		TLSConfig: tlsConfig,
		Handler:   mux,
	})
	err = svr.Start()
	if err != nil {
		fmt.Printf("nbio.Start failed: %v\n", err)
		return
	}
	defer svr.Stop()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
	log.Println("exit")
}
