package middlewares

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

var COOKIES_FIELDS = []string{"auth"}

func CheackAdminIsValid(c *fiber.Ctx) error {
	for i := 0; i < len(COOKIES_FIELDS); i++ {
		cookie := c.Cookies(COOKIES_FIELDS[i])

		if cookie != os.Getenv("ADMIN_TOKEN") {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}
	}
	return c.Next()
}
