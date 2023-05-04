package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func Setup(router fiber.Router) {
	router.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	router.Use(requestid.New())
	router.Use(logger.New(logger.Config{
		Format:     "[${time}] ${pid} [${ip}]:${port} ${locals:requestid} ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "2006-01-02T15:04:05 MST",
	}))
}
