package utils

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var Validate = validator.New()

func ParseJSON(c *fiber.Ctx, payload any) error {
	if c.Body() == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Missing Request Body")
	}

	err := json.NewDecoder(bytes.NewReader(c.Body())).Decode(&payload)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid payload")
	}

	return nil
}

func ParseAndValidateJSON(c *fiber.Ctx, v interface{}) error {
	if err := ParseJSON(c, v); err != nil {
		return err
	}
	if err := Validate.Struct(v); err != nil {
		return fmt.Errorf("invalid payload: %v", err)
	}
	return nil
}
