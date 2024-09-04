package user

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/zondaf12/planner-app-backend/config"
	"github.com/zondaf12/planner-app-backend/service/auth"
	"github.com/zondaf12/planner-app-backend/types"
	"github.com/zondaf12/planner-app-backend/utils"
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
	// Get JSON Payload
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(c, &payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("invalid payload %v", errors).Error())
	}

	// Check if the user already exists
	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email).Error())
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// If it doesnt exist, create a new user
	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *Handler) HandleLogin(c *fiber.Ctx) error {
	// Parse payload
	var payload types.LoginUserPayload
	if err := utils.ParseJSON(c, &payload); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	// Validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		return fiber.NewError(http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors).Error())
	}

	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, fmt.Errorf("not found, invalid email or password").Error())
	}

	if !auth.ComparePassword(u.Password, payload.Password) {
		return fiber.NewError(http.StatusBadRequest, fmt.Errorf("not found, invalid email or password").Error())
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	c.Response().Header.Set("Access-Control-Expose-Headers", "X-Token")
	c.Response().Header.Set("X-Token", token)

	return c.Status(http.StatusOK).JSON(fiber.Map{"userId": u.ID.String()})
}
