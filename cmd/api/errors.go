package main

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (app *Application) internalServerError(c *fiber.Ctx, err error) error {
	log.Printf("internal server error: %s path: %s error: %v", c.Method(), c.Path(), err.Error())

	return writeJSONError(c, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *Application) badRequestResponse(c *fiber.Ctx, err error) error {
	log.Printf("bad request error: %s path: %s error: %v", c.Method(), c.Path(), err.Error())

	return writeJSONError(c, http.StatusBadRequest, err.Error())
}

func (app *Application) conflictResponse(c *fiber.Ctx, err error) error {
	log.Printf("conflict error: %s path: %s error: %v", c.Method(), c.Path(), err.Error())

	return writeJSONError(c, http.StatusConflict, err.Error())
}

func (app *Application) notFoundResponse(c *fiber.Ctx, err error) error {
	log.Printf("not found error: %s path: %s error: %v", c.Method(), c.Path(), err.Error())

	return writeJSONError(c, http.StatusNotFound, "not found")
}
