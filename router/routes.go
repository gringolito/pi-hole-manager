package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gringolito/pi-hole-manager/host"
)

func SetupRoutes(router fiber.Router, staticHostsFilePath string) {
	api := router.Group("/api")
	v1 := api.Group("/v1")

	hostController := host.NewController(staticHostsFilePath)
	hostController.SetupRoutes(v1)
}
