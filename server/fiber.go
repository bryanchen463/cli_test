package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

var messages map[int][]byte

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var rd = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateRandomString(length int) string {

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rd.Intn(len(charset))]
	}

	return string(b)
}

func init() {
	messages = make(map[int][]byte)
	counts := []int{100, 200, 500, 1024}
	for _, c := range counts {
		m := generateRandomString(c)
		messages[c] = []byte(m)
	}
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})
	app.Get("/echo", func(c *fiber.Ctx) error {
		queryParam1 := c.FormValue("name", "default_value")
		return c.SendString(queryParam1)
	})
	app.Post("/echo", func(c *fiber.Ctx) error {
		// 获取查询参数
		queryParam1 := c.FormValue("name", "default_value")
		return c.SendString(queryParam1)
	})

	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/echo", websocket.New(func(c *websocket.Conn) {
		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				break
			}

			if err = c.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				break
			}
		}

	}))

	app.Get("/ws/message/:size/:count", websocket.New(func(c *websocket.Conn) {
		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		size := c.Params("size")
		count := c.Params("count")
		s, err := strconv.Atoi(size)
		if err != nil {
			c.WriteMessage(websocket.CloseUnsupportedData, []byte("size param invalid"))
			return
		}
		n, err := strconv.Atoi(count)
		if err != nil {
			c.WriteMessage(websocket.CloseUnsupportedData, []byte("count param invalid"))
			return
		}
		message, ok := messages[s]
		if !ok {
			log.Println("unkown size", s)
			c.WriteMessage(websocket.CloseAbnormalClosure, []byte("unkown size"))
			return
		}
		for i := 0; i < n; i++ {
			if err = c.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("write:", err)
				return
			}
		}
		c.WriteMessage(websocket.TextMessage, []byte("finish"))
	}))

	go func() {
		quickwsServer()
	}()

	// go func() {
	// 	nbioWsServer()
	// }()

	// go func() {
	// 	fastHttpWs()
	// }()

	err := app.ListenTLS(":8090", "../nginx/cert/server.crt", "../nginx/cert/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
