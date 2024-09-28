package main

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (app *Application) BasicAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// read the username and password from the request
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return app.unauthorizedBasicErrorResponse(c, fmt.Errorf("no authorization header"))
		}

		// parse the username and password from the request
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Basic" {
			return app.unauthorizedBasicErrorResponse(c, fmt.Errorf("invalid authorization header"))
		}

		// decode the username and password
		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			return app.unauthorizedBasicErrorResponse(c, err)
		}

		// check if the username and password are correct
		username := app.config.auth.basic.username
		password := app.config.auth.basic.password

		creds := strings.SplitN(string(decoded), ":", 2)
		if len(creds) != 2 || creds[0] != username || creds[1] != password {
			return app.unauthorizedBasicErrorResponse(c, fmt.Errorf("invalid credentials"))
		}

		return c.Next()
	}

}
