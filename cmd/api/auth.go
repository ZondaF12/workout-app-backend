package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zondaf12/workout-app-backend/internal/mailer"
	"github.com/zondaf12/workout-app-backend/internal/store"
)

type RegisterUserPayload struct {
	Username  string `json:"username" validate:"required,max=100"`
	Email     string `json:"email" validate:"required,email,max=255"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required,max=100"`
	LastName  string `json:"last_name" validate:"required,max=100"`
}

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	UserWithToken		"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/register [post]
func (app *Application) registerUserHandler(c *fiber.Ctx) error {
	var payload RegisterUserPayload
	if err := readJSON(c, &payload); err != nil {
		return app.badRequestResponse(c, err)
	}

	if err := Validate.Struct(payload); err != nil {
		return app.badRequestResponse(c, err)
	}

	user := store.User{
		Username:  payload.Username,
		Email:     payload.Email,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
	}

	// Hash the password
	if err := user.Password.Set(payload.Password); err != nil {
		return app.internalServerError(c, err)
	}

	plainToken := uuid.New().String()

	hash := sha256.Sum256([]byte(plainToken))
	hashedToken := hex.EncodeToString(hash[:])

	// Store the user
	err := app.store.Users.CreateAndInvite(c.Context(), &user, hashedToken, app.config.mail.exp)
	if err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			return app.badRequestResponse(c, err)
		case store.ErrDuplicateUsername:
			return app.badRequestResponse(c, err)
		default:
			return app.internalServerError(c, err)
		}
	}

	userWithToken := UserWithToken{
		User:  &user,
		Token: plainToken,
	}

	isProdEnv := app.config.env == "production"

	// TODO: Change this to the frontend app url
	activationUrl := fmt.Sprintf("%s/v1/authentication/activate?token=%s", app.config.apiUrl, plainToken)
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationUrl,
	}

	// Send the email
	status, err := app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {
		app.logger.Errorw("error sending welcome email", "error", err)

		// rollback user creation if email fails
		if err := app.store.Users.Delete(c.Context(), user.ID); err != nil {
			app.logger.Errorw("error deleting user", "error", err)
		}

		return app.internalServerError(c, err)
	}

	app.logger.Infow("Email sent", "status code", status)

	if err := app.jsonResponse(c, http.StatusCreated, userWithToken); err != nil {
		return app.internalServerError(c, err)
	}

	return nil
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8"`
}

// createTokenHandler godoc
//
//	@Summary		Creates a token
//	@Description	Creates a token for a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	string	"Token"
//	@Failure		400	{object}	error
//	@Failure		401	{object}	error
//	@Failure		500	{object}	error
//	@Router			/authentication/login [post]
func (app *Application) createTokenHandler(c *fiber.Ctx) error {
	self := getSelfFromContext(c)

	// generate a token for the user & add claims
	claims := jwt.MapClaims{
		"sub": self.ID,
		"exp": time.Now().Add(app.config.auth.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.token.iss,
		"aud": app.config.auth.token.iss,
	}

	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		return app.internalServerError(c, err)
	}

	// set the token in the cookie
	setTokenCookie(c, token)

	if err := app.jsonResponse(c, http.StatusCreated, self); err != nil {
		return app.internalServerError(c, err)
	}

	return nil
}

func setTokenCookie(c *fiber.Ctx, token string) {
	cookie := new(fiber.Cookie)
	cookie.Name = "jwt"
	cookie.Value = token
	cookie.Expires = time.Now().Add(time.Hour * 24 * 3) // Set expiration time
	cookie.HTTPOnly = true                              // Makes the cookie inaccessible to client-side scripts
	cookie.Secure = true                                // Only send over HTTPS
	cookie.SameSite = "Lax"                             // Provides some CSRF protection

	c.Cookie(cookie)
}
