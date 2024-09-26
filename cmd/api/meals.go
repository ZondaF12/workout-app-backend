package main

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/zondaf12/workout-app-backend/internal/store"
)

type CreateMealEntryPayload struct {
	FoodID      string  `json:"food_id" validate:"required"`
	ServingUnit string  `json:"serving_unit" validate:"required"`
	Amount      float64 `json:"amount" validate:"required"`
	ConsumedAt  string  `json:"consumed_at" validate:"required"`
	MealName    string  `json:"meal_name" validate:"required"`
}

type UpdateMealEntryPayload struct {
	ServingUnit string  `json:"serving_unit" validate:"required"`
	Amount      float64 `json:"amount" validate:"required"`
	ConsumedAt  string  `json:"consumed_at" validate:"required"`
	MealName    string  `json:"meal_name" validate:"required"`
}

// CreateMealEntry godoc
//
//	@Summary		Creates a meal entry
//	@Description	Creates a meal entry
//	@Tags			meal entrys
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateMealEntryPayload	true	"Meal payload"
//	@Success		201		{object}	store.MealEntry
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/meal [post]
func (app *Application) createMealEntryHandler(c *fiber.Ctx) error {
	var payload CreateMealEntryPayload
	if err := readJSON(c, &payload); err != nil {
		return app.badRequestResponse(c, err)
	}

	if err := Validate.Struct(payload); err != nil {
		return app.badRequestResponse(c, err)
	}

	checkMeal := store.Meal{
		UserID: uuid.MustParse("9ecafdec-7cc0-451a-847d-2aa02cf2adc5"),
		Name:   payload.MealName,
		Date:   payload.ConsumedAt,
	}

	meal, err := app.store.Meals.GetMeal(c.Context(), checkMeal)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			if err := app.store.Meals.CreateMeal(c.Context(), &checkMeal); err != nil {
				return app.internalServerError(c, err)
			}
		default:
			return app.internalServerError(c, err)
		}
	}

	var mealID uuid.UUID
	if meal != nil {
		mealID = meal.ID
	} else {
		mealID = checkMeal.ID
	}

	newEntry := store.MealEntry{
		MealID:      mealID,
		FoodID:      uuid.MustParse(payload.FoodID),
		ServingUnit: payload.ServingUnit,
		Amount:      payload.Amount,
		ConsumedAt:  payload.ConsumedAt,
	}
	if err := app.store.Meals.CreateMealEntry(c.Context(), &newEntry); err != nil {
		return app.internalServerError(c, err)
	}

	if err := app.jsonResponse(c, fiber.StatusCreated, newEntry); err != nil {
		return app.internalServerError(c, err)
	}

	return nil
}

// UpdateMealEntry godoc
//
//	@Summary		Updates a meal entry
//	@Description	Updates a meal entry by ID
//	@Tags			meal entrys
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Meal ID"
//	@Param			payload	body		UpdateMealEntryPayload	true	"Meal payload"
//	@Success		200		{object}	store.Meal
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/meal/{id} [patch]
func (app *Application) updateMealEntryHandler(c *fiber.Ctx) error {
	mealID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return app.badRequestResponse(c, err)
	}

	var payload UpdateMealEntryPayload
	if err := readJSON(c, &payload); err != nil {
		return app.badRequestResponse(c, err)
	}

	if err := Validate.Struct(payload); err != nil {
		return app.badRequestResponse(c, err)
	}

	/* Get the meal entry */
	entry, err := app.store.Meals.GetMealEntryByID(c.Context(), mealID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			return app.notFoundResponse(c, err)
		default:
			return app.internalServerError(c, err)
		}
	}

	/* Get the current meal from that entry */
	currentMeal, err := app.store.Meals.GetMealByID(c.Context(), entry.MealID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			return app.notFoundResponse(c, err)
		default:
			return app.internalServerError(c, err)
		}
	}

	/* Compare to see if the payload changes the meal that the entry is on */
	updatedEntry := store.MealEntry{
		ID:          mealID,
		MealID:      entry.MealID,
		FoodID:      entry.FoodID,
		ServingUnit: payload.ServingUnit,
		Amount:      payload.Amount,
		ConsumedAt:  payload.ConsumedAt,
	}
	if currentMeal.Name != payload.MealName {
		checkMeal := store.Meal{
			UserID: uuid.MustParse("9ecafdec-7cc0-451a-847d-2aa02cf2adc5"),
			Name:   payload.MealName,
			Date:   payload.ConsumedAt,
		}

		meal, err := app.store.Meals.GetMeal(c.Context(), checkMeal)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				if err := app.store.Meals.CreateMeal(c.Context(), &checkMeal); err != nil {
					return app.internalServerError(c, err)
				}
			default:
				return app.internalServerError(c, err)
			}
		}

		var mealID uuid.UUID
		if meal != nil {
			mealID = meal.ID
		} else {
			mealID = checkMeal.ID
		}

		updatedEntry.MealID = mealID
	}

	if err := app.store.Meals.UpdateMealEntry(c.Context(), &updatedEntry); err != nil {
		return app.internalServerError(c, err)
	}

	if err := app.jsonResponse(c, fiber.StatusOK, updatedEntry); err != nil {
		return app.internalServerError(c, err)
	}

	return nil
}

// DeleteMealEntry godoc
//
//	@Summary		Deletes a meal entry
//	@Description	Delete a meal entry by ID
//	@Tags			meal entrys
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Meal ID"
//	@Success		204	{object} string
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/meal/{id} [delete]
func (app *Application) deleteMealEntryHandler(c *fiber.Ctx) error {
	mealID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return app.badRequestResponse(c, err)
	}

	if err := app.store.Meals.DeleteMealEntry(c.Context(), mealID); err != nil {
		return app.internalServerError(c, err)
	}

	if err := app.jsonResponse(c, http.StatusNoContent, nil); err != nil {
		return app.internalServerError(c, err)
	}

	return nil
}
