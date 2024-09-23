package main

import (
	"github.com/gofiber/fiber/v2"
)

func (app *Application) healthCheckHandler(c *fiber.Ctx) error {
	return c.SendString("OK")
}
