package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gringolito/pi-hole-manager/api/middleware"
	"github.com/gringolito/pi-hole-manager/api/routes"
	"github.com/gringolito/pi-hole-manager/pkg/host"
)

const STATIC_HOSTS_FILE_PATH string = "./04-pihole-static-dhcp.conf"

func addHostApi(router fiber.Router) {
	hostRepository := host.NewRepository(STATIC_HOSTS_FILE_PATH)
	hostService := host.NewService(hostRepository)
	routes.HostRouter(router, hostService)
}

func main() {
	app := fiber.New(fiber.Config{
		CaseSensitive:     true,
		EnablePrintRoutes: true,
		AppName:           "Pi-Hole Manager v0.1.0",
	})

	middleware.Setup(app)

	routes.MetricsRouter(app)

	api := app.Group("/api")
	v1 := api.Group("/v1")
	addHostApi(v1)

	log.Fatal(app.Listen(":8080"))
}
