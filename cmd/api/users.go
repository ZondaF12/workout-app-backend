package main

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/zondaf12/workout-app-backend/internal/store"
)

type selfKey string

const selfCtxKey selfKey = "self"

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//
//	@Security		ApiKeyAuth
//
//	@Router			/users/{id} [get]
func (app *Application) getUserHandler(c *fiber.Ctx) error {
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
//	@Param			userID	path		string	true	"User ID"
//	@Success		204		{string}	string	"User followed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//
//	@Security		ApiKeyAuth
//
//	@Router			/users/{userID}/follow [put]
func (app *Application) followUserHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	followedUserID, err := uuid.Parse(id)
	if err != nil {
		return app.internalServerError(c, err)
	}

	self := getSelfFromContext(c)

	if err := app.store.Followers.Follow(c.Context(), followedUserID, self.ID); err != nil {
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
//	@Param			userID	path		string	true	"User ID"
//	@Success		204		{string}	string	"User unfollowed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//
//	@Security		ApiKeyAuth
//
//	@Router			/users/{userID}/unfollow [put]
func (app *Application) unfollowUserHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	followedUserID, err := uuid.Parse(id)
	if err != nil {
		return app.internalServerError(c, err)
	}

	self := getSelfFromContext(c)

	if err := app.store.Followers.Unfollow(c.Context(), followedUserID, self.ID); err != nil {
		return app.internalServerError(c, err)
	}

	if err := app.jsonResponse(c, http.StatusNoContent, nil); err != nil {
		return app.internalServerError(c, err)
	}

	return nil
}

// ActivateUser godoc
//
//	@Summary		Activates/Register a user
//	@Description	Activates/Register a user by invitation token
//	@Tags			users
//	@Produce		json
//	@Param			token	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"User activated"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put]
func (app *Application) activateUserHandler(c *fiber.Ctx) error {
	token := c.Params("token")

	if err := app.store.Users.Activate(c.Context(), token); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			return app.notFoundResponse(c, err)
		default:
			return app.internalServerError(c, err)
		}
	}

	if err := app.jsonResponse(c, http.StatusNoContent, nil); err != nil {
		return app.internalServerError(c, err)
	}

	return nil
}

func getSelfFromContext(c *fiber.Ctx) *store.User {
	self, _ := c.Locals(selfCtxKey).(*store.User)

	return self
}
