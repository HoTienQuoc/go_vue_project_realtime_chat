package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(c fiber.Ctx) error {
	tok := c.Get("Authorization")
	if tok == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}
	// check if the token starts with "Bearer"
	if !strings.HasPrefix(tok, "Bearer ") {
		tok = strings.TrimSpace(tok)
	} else {
		splited := strings.Split(tok, "Bearer ")
		if len(splited) != 2 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "unauthenticated",
			})
		}
		tok = splited[1]
	}

	SecretKey := os.Getenv("JWT_SECRET")
	token, err := jwt.ParseWithClaims(tok, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)

	if !ok || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	c.Locals("userId", claims.Issuer)
	return c.Next()
}
