package api

import (
	"github.com/gofiber/fiber/v2"

	"github.com/phgh1246/golang_project01/types"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return ErrorUnauthorized()
	}
	if !user.IsAdmin {
		return ErrorUnauthorized()
	}
	return c.Next()
}
