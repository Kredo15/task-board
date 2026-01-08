package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type JWTMiddleware struct {
	secretKey string
}

func NewJWTMiddleware(s string) *JWTMiddleware {
	return &JWTMiddleware{
		secretKey: s,
	}
}

// AuthRequired проверяет наличие и валидность JWT токена
func (m *JWTMiddleware) VerifyToken(c *fiber.Ctx) error {
	// Извлекаем заголовок Authorization
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing authorization header",
		})
	}

	// Проверяем формат "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid authorization format",
		})
	}

	tokenString := parts[1]

	// Парсим и валидируем токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Проверка алгоритма подписи (защита от атак подмены алгоритма)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Unexpected signing method")
		}
		return []byte(m.secretKey), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	// 4. Извлекаем данные (claims) и сохраняем в контекст для следующих обработчиков
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		// Например, сохраняем user_id
		c.Locals("user_id", claims["sub"])
	}

	return c.Next()
}
