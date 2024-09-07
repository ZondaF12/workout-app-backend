package user

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/zondaf12/workout-app-backend/config"
	"github.com/zondaf12/workout-app-backend/service/auth"
	"github.com/zondaf12/workout-app-backend/types"
	"github.com/zondaf12/workout-app-backend/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Post("/register", h.HandleRegister)
	router.Post("/login", h.HandleLogin)
}

func (h *Handler) HandleRegister(c *fiber.Ctx) error {
	var payload types.RegisterUserPayload
	if err := utils.ParseAndValidateJSON(c, &payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Check if the user already exists
	if _, err := h.store.GetUserByEmail(payload.Email); err == nil {
		return fiber.NewError(fiber.StatusConflict, "user with this email already exists")
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	// Create a new user
	user := types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	}
	if err := h.store.CreateUser(user); err != nil {
		log.Printf("Error creating user: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *Handler) HandleLogin(c *fiber.Ctx) error {
	// Parse payload
	var payload types.LoginUserPayload
	if err := utils.ParseAndValidateJSON(c, &payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, fmt.Errorf("not found, invalid email or password").Error())
	}

	if !auth.ComparePassword(u.Password, payload.Password) {
		return fiber.NewError(http.StatusBadRequest, fmt.Errorf("not found, invalid email or password").Error())
	}

	token, err := auth.CreateJWT([]byte(config.Envs.JWTSecret), u.ID)
	if err != nil {
		log.Printf("Error creating JWT: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	c.Set("Access-Control-Expose-Headers", "X-Token")
	c.Set("X-Token", token)

	return c.Status(http.StatusOK).JSON(fiber.Map{"userId": u.ID.String()})
}
