package user

import "github.com/gofiber/fiber/v2"

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Post("/register", h.HandleRegister)
	router.Post("/login", h.HandleLogin)
}

func (h *Handler) HandleRegister(c *fiber.Ctx) error {
	return c.SendString("Register")
}

func (h *Handler) HandleLogin(c *fiber.Ctx) error {
	return c.SendString("Login")
}
