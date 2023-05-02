package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gringolito/pi-hole-manager/router"
)

const STATIC_HOSTS_FILE_PATH string = "./04-pihole-static-dhcp.conf"

func main() {
	app := fiber.New(fiber.Config{
		// Prefork:       true,
		// StrictRouting:     true,
		CaseSensitive:     true,
		EnablePrintRoutes: true,
		AppName:           "Pi-Hole Manager v0.1.0",
	})

	router.SetupRoutes(app, STATIC_HOSTS_FILE_PATH)

	log.Fatal(app.Listen(":8080"))
}
