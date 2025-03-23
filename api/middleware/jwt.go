package middleware

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("-- JWT auth")

	tokenString, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return fmt.Errorf("unauthorized")
	}

	if err := parseToken(tokenString[0]); err != nil {
		return err
	}

	return nil
}

func parseToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		fmt.Println("GETENV_TOKEN---->>>", secret)
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse JWT token:", err)
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
	}

	return fmt.Errorf("unauthorized")
}
