package main

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (app *Application) AuthTokenMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// read the token from the request
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return app.unauthorizedErrorResponse(c, fmt.Errorf("no authorization header"))
		}

		// parse the token from the request
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return app.unauthorizedErrorResponse(c, fmt.Errorf("invalid authorization header"))
		}

		// validate the token
		token := parts[1]
		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			return app.unauthorizedErrorResponse(c, err)
		}

		// extract the user ID from the token
		claims := jwtToken.Claims.(jwt.MapClaims)

		userID, err := uuid.Parse(claims["sub"].(string))
		if err != nil {
			return app.unauthorizedErrorResponse(c, err)
		}

		// get the user from the store
		user, err := app.store.Users.GetByID(c.Context(), userID)
		if err != nil {
			return app.unauthorizedErrorResponse(c, err)
		}

		// set the user in the context
		c.Locals(selfCtxKey, user)

		return c.Next()
	}
}

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
