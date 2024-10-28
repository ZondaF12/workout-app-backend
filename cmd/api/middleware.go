package main

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zondaf12/workout-app-backend/internal/store"
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
		user, err := app.getUser(c, userID)
		if err != nil {
			return app.unauthorizedErrorResponse(c, err)
		}

		// set the user in the context
		c.Locals(selfCtxKey, user)

		return c.Next()
	}
}

func (app *Application) getUser(c *fiber.Ctx, userID uuid.UUID) (*store.User, error) {
	if !app.config.redisCfg.enabled {
		return app.store.Users.GetByID(c.Context(), userID)
	}

	user, err := app.cacheStorage.Users.Get(c.Context(), userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, err = app.store.Users.GetByID(c.Context(), userID)
		if err != nil {
			return nil, err
		}

		if err := app.cacheStorage.Users.Set(c.Context(), user); err != nil {
			return nil, err
		}
	}

	return user, nil
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

		creds := strings.SplitN(string(decoded), ":", 2)

		// fetch the user from the payload (check if the user exists)
		user, err := app.store.Users.GetByEmail(c.Context(), creds[0])
		if err != nil {
			switch err {
			case store.ErrNotFound:
				return app.unauthorizedErrorResponse(c, err)
			default:
				return app.internalServerError(c, err)
			}
		}

		if !user.Password.Verify(creds[1]) {
			return app.unauthorizedBasicErrorResponse(c, fmt.Errorf("invalid credentials"))
		}

		// set the user in the context
		c.Locals(selfCtxKey, user)

		return c.Next()
	}
}

func (app *Application) checkMealOwnership() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := getSelfFromContext(c)
		meal := getMealEntryFromContext(c)

		if meal.Meal.UserID == user.ID {
			return c.Next()
		}

		return app.forbiddenResponse(c)
	}
}

func (app *Application) RateLimiterMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if app.config.rateLimiter.Enabled {
			if allow, retryAfter := app.rateLimiter.Allow(c.IP()); !allow {
				return app.rateLimitExceededResponse(c, retryAfter.String())
			}
		}
		return c.Next()
	}
}
