package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gringolito/pi-hole-manager/router"
)

const STATIC_HOSTS_FILE_PATH string = "./04-pihole-static-dhcp.conf"

func main() {
	app := fiber.New(fiber.Config{
		// Prefork:           true,
		// StrictRouting:     true,
		CaseSensitive:     true,
		EnablePrintRoutes: true,
		AppName:           "Pi-Hole Manager v0.1.0",
	})

	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "${pid} ${locals:requestid} [${ip}]:${port} ${status} - ${method} ${path}​\n",
	}))

	router.SetupRoutes(app, STATIC_HOSTS_FILE_PATH)

	log.Fatal(app.Listen(":8080"))
}
