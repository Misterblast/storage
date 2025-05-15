package utils

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Cors() fiber.Handler {
	origins := strings.Split(os.Getenv("CORS_ORIGIN"), ",")
	methods := os.Getenv("CORS_METHODS")

	return cors.New(
		cors.Config{
			AllowOrigins: strings.Join(origins, ","),
			AllowMethods: methods,
		},
	)
}
