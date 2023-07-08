package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gringolito/pi-hole-manager/api/handler"
	"github.com/gringolito/pi-hole-manager/api/middleware/fiberswagger"
	"github.com/gringolito/pi-hole-manager/pkg/host"
)

func HostRouter(api fiber.Router, service host.Service) {
	api.Route("/static", func(router fiber.Router) {
		router.Get("/hosts", handler.GetAllStaticHosts(service)).Name("get_all")
		router.Get("/host", handler.GetStaticHost(service)).Name("get")
		router.Post("/host", handler.AddStaticHost(service)).Name("add")
		router.Put("/host", handler.UpdateStaticHost(service)).Name("update")
		router.Delete("/host", handler.RemoveStaticHost(service)).Name("remove")
	}, "static.hosts.")
}

func MetricsRouter(router fiber.Router, cfg monitor.Config) {
	router.Get("/metrics", monitor.New(cfg))
}

func OpenApiRouter(router fiber.Router, cfg fiberswagger.Config) {
	cfg.BasePath = "/openapi"
	router.Use("/openapi", fiberswagger.New(cfg))
}
