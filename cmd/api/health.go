package main

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// healthcheckHandler godoc
//
//	@Summary		Healthcheck
//	@Description	Healthcheck endpoint
//	@Tags			tools
//	@Produce		json
//	@Success		200	{object}	string	"ok"
//	@Router			/health [get]
func (app *Application) healthCheckHandler(c *fiber.Ctx) error {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}

	if err := app.jsonResponse(c, http.StatusOK, data); err != nil {
		return app.internalServerError(c, err)
	}

	return nil
}
