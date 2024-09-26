package main

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/zondaf12/workout-app-backend/internal/store"
)

type userKey string

const userCtxKey userKey = "user"

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//
//	@Security		ApiKeyAuth
//
//	@Router			/user/{id} [get]
func (app *Application) getUserHandler(c *fiber.Ctx) error {
	user := getUserFromContext(c)

	if err := app.jsonResponse(c, http.StatusOK, user); err != nil {
		return app.internalServerError(c, err)
	}

	return nil
}

// FollowUser godoc
//
//	@Summary		Follows a user
//	@Description	Follows a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User followed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//
//	@Security		ApiKeyAuth
//
//	@Router			/users/{userID}/follow [put]
func (app *Application) followUserHandler(c *fiber.Ctx) error {
	followedUser := getUserFromContext(c)

	// TODO: Get the user ID from auth
	userId := uuid.MustParse("9ecafdec-7cc0-451a-847d-2aa02cf2adc5")

	if err := app.store.Followers.Follow(c.Context(), followedUser.ID, userId); err != nil {
		switch err {
		case store.ErrConflict:
			return app.conflictResponse(c, err)
		default:
			return app.internalServerError(c, err)
		}
	}

	if err := app.jsonResponse(c, http.StatusNoContent, nil); err != nil {
		return app.internalServerError(c, err)
	}

	return nil
}

// UnfollowUser gdoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User unfollowed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//
//	@Security		ApiKeyAuth
//
//	@Router			/users/{userID}/unfollow [put]
func (app *Application) unfollowUserHandler(c *fiber.Ctx) error {
	unfollowedUser := getUserFromContext(c)

	// TODO: Get the user ID from auth
	userId := uuid.MustParse("9ecafdec-7cc0-451a-847d-2aa02cf2adc5")

	if err := app.store.Followers.Unfollow(c.Context(), unfollowedUser.ID, userId); err != nil {
		return app.internalServerError(c, err)
	}

	if err := app.jsonResponse(c, http.StatusNoContent, nil); err != nil {
		return app.internalServerError(c, err)
	}

	return nil
}

func (app *Application) userContextMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			return app.internalServerError(c, err)
		}

		user, err := app.store.Users.GetByID(c.Context(), uuid)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				return app.notFoundResponse(c, err)
			default:
				return app.internalServerError(c, err)
			}
		}

		c.Locals(userCtxKey, user)
		return c.Next()
	}
}

func getUserFromContext(c *fiber.Ctx) *store.User {
	user, _ := c.Locals(userCtxKey).(*store.User)

	return user
}
