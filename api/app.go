package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gringolito/pi-hole-manager/api/middleware"
	"github.com/gringolito/pi-hole-manager/api/routes"
	"github.com/gringolito/pi-hole-manager/config"
	"github.com/gringolito/pi-hole-manager/pkg/host"
)

const (
	AppName    = "Pi-Hole Manager"
	AppVersion = "v0.1.0"
)

func addHostApi(router fiber.Router, cfg *config.Config) {
	hostRepository := host.NewRepository(cfg.Host.Static.File)
	hostService := host.NewService(hostRepository)
	routes.HostRouter(router, hostService)
}

func main() {
	cfg, err := config.Init("test")
	if err != nil {
		log.Fatal(err.Error())
	}

	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		// EnablePrintRoutes: true,
		AppName: fmt.Sprintf("%s %s", AppName, AppVersion),
	})

	middleware.Setup(app)

	routes.MetricsRouter(app, monitor.Config{
		Title: fmt.Sprintf("%s Monitor", AppName),
	})

	api := app.Group("/api")
	v1 := api.Group("/v1")
	addHostApi(v1, cfg)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", cfg.Server.Port)))
}
