package main

import (
	"github.com/gofiber/fiber/v2"
)

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

	app.Listen(":8090")
}
