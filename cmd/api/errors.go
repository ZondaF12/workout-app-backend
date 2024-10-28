package main

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (app *Application) internalServerError(c *fiber.Ctx, err error) error {
	app.logger.Errorw("internal server error", "method", c.Method(), "path", c.Path(), "error", err.Error())

	return writeJSONError(c, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *Application) badRequestResponse(c *fiber.Ctx, err error) error {
	app.logger.Warnf("bad request error", "method", c.Method(), "path", c.Path(), "error", err.Error())

	return writeJSONError(c, http.StatusBadRequest, err.Error())
}

func (app *Application) conflictResponse(c *fiber.Ctx, err error) error {
	app.logger.Errorw("conflict error", "method", c.Method(), "path", c.Path(), "error", err.Error())

	return writeJSONError(c, http.StatusConflict, err.Error())
}

func (app *Application) notFoundResponse(c *fiber.Ctx, err error) error {
	app.logger.Warnf("not found error", "method", c.Method(), "path", c.Path(), "error", err.Error())

	return writeJSONError(c, http.StatusNotFound, "not found")
}

func (app *Application) unauthorizedErrorResponse(c *fiber.Ctx, err error) error {
	app.logger.Warnf("unauthorized error", "method", c.Method(), "path", c.Path(), "error", err.Error())

	return writeJSONError(c, http.StatusUnauthorized, "unauthorized")
}

func (app *Application) unauthorizedBasicErrorResponse(c *fiber.Ctx, err error) error {
	app.logger.Warnf("unauthorized basic error", "method", c.Method(), "path", c.Path(), "error", err.Error())

	c.Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	return writeJSONError(c, http.StatusUnauthorized, "unauthorized")
}

func (app *Application) forbiddenResponse(c *fiber.Ctx) error {
	app.logger.Warnf("forbidden", "method", c.Method(), "path", c.Path(), "error")

	return writeJSONError(c, http.StatusForbidden, "forbidden")
}

func (app *Application) rateLimitExceededResponse(c *fiber.Ctx, retryAfter string) error {
	app.logger.Warnw("rate limit exceeded", "method", c.Method(), "path", c.Path())

	c.Set("Retry-After", retryAfter)

	return writeJSONError(c, http.StatusTooManyRequests, "rate limit exceeded, retry after:"+retryAfter)
}
