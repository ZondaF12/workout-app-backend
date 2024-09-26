package main

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/zondaf12/workout-app-backend/internal/store"
)

// TODO: Update validate struct tags
type CreateFoodPayload struct {
	Name        string  `json:"name" validate:"required,max=100"`
	Description string  `json:"description" validate:"required,max=1000"`
	Calories    int     `json:"calories" validate:"required"`
	Protein     float64 `json:"protein" validate:"required"`
	Carbs       float64 `json:"carbs" validate:"required"`
	Fat         float64 `json:"fat" validate:"required"`
	Brand       string  `json:"brand" validate:"max=100"`
	ServingSize float64 `json:"serving_size" validate:"required"`
	ServingUnit string  `json:"serving_unit" validate:"required"`
	Verified    bool    `json:"verified"`
}

// CreateFood godoc
//
//	@Summary		Creates a food
//	@Description	Creates a food
//	@Tags			foods
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateFoodPayload	true	"Food payload"
//	@Success		201		{object}	store.Food
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/food [post]
func (app *Application) createFoodHandler(c *fiber.Ctx) error {
	var payload CreateFoodPayload
	if err := readJSON(c, &payload); err != nil {
		return app.badRequestResponse(c, err)
	}

	if err := Validate.Struct(payload); err != nil {
		return app.badRequestResponse(c, err)
	}

	food := store.Food{
		Name:        payload.Name,
		Description: payload.Description,
		Calories:    payload.Calories,
		Protein:     payload.Protein,
		Carbs:       payload.Carbs,
		Fat:         payload.Fat,
		Brand:       payload.Brand,
		ServingSize: payload.ServingSize,
		ServingUnit: payload.ServingUnit,
	}

	if err := app.store.Foods.Create(c.Context(), &food); err != nil {
		return app.internalServerError(c, err)
	}

	if err := app.jsonResponse(c, http.StatusCreated, food); err != nil {
		return app.internalServerError(c, err)
	}

	return nil
}

// GetFoo godoc
//
//	@Summary		Fetches a food
//	@Description	Fetches a food by ID
//	@Tags			foods
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Food ID"
//	@Success		200	{object}	store.Food
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/food/{id} [get]
func (app *Application) getFoodHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return app.internalServerError(c, err)
	}

	food, err := app.store.Foods.GetByID(c.Context(), uuid)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			return app.notFoundResponse(c, err)
		default:
			return app.internalServerError(c, err)
		}
	}

	if err := app.jsonResponse(c, http.StatusOK, food); err != nil {
		return app.internalServerError(c, err)
	}

	return nil
}
