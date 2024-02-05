package middlewares

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func CheckAdminIsValid(c *fiber.Ctx) error {
	// Получаем значение токена из заголовка 'Authorization'
	token := c.Get("Authorization")

	// Проверяем, что токен присутствует и соответствует ожидаемому значению
	if token == "" || token != "Bearer "+os.Getenv("ADMIN_TOKEN") {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	// Если все проверки пройдены, переходим к следующему middleware или обработчику
	return c.Next()
}
