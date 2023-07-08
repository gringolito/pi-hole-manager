package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gringolito/pi-hole-manager/api/middleware/fiberslog"
	"golang.org/x/exp/slog"
)

func Setup(router fiber.Router, logger *slog.Logger) {
	router.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	router.Use(requestid.New())

	router.Use(fiberslog.New(fiberslog.Config{
		Logger: logger,
		Fields: []string{"latency", "status", "method", "path", "requestId", "ip", "port", "pid"},
	}))
}
