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
