package foods

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zondaf12/workout-app-backend/types"
)

type Handler struct {
	store types.FoodStore
}

func NewHandler(store types.FoodStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Post("/food", h.CreateFood)
	router.Get("/food", h.GetAllFoods)
}

func (h *Handler) CreateFood(c *fiber.Ctx) error {
	return nil
}

func (h *Handler) GetAllFoods(c *fiber.Ctx) error {
	foods, err := h.store.GetAllFoods()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(foods)
}
