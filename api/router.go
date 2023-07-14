package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gringolito/pi-hole-manager/api/handler"
	"github.com/gringolito/pi-hole-manager/api/middleware/fiberswagger"
	"github.com/gringolito/pi-hole-manager/api/scope"
	"github.com/gringolito/pi-hole-manager/pkg/host"
)

type Router struct {
	root  fiber.Router
	api   fiber.Router
	apiv1 fiber.Router
	mw    Middleware
}

func NewRouter(root fiber.Router, mw Middleware) Router {
	root.Use(mw.Recovery())
	root.Use(mw.Logger())

	api := root.Group("/api")
	api.Use(mw.RequestId())

	apiv1 := api.Group("/v1")

	return Router{
		root:  root,
		api:   api,
		apiv1: apiv1,
		mw:    mw,
	}
}

func (r Router) HostApi(service host.Service) {
	r.apiv1.Route("/static", func(router fiber.Router) {
		router.Get("/hosts", r.mw.Authentication(scope.DhcpCanRead...), handler.GetAllStaticHosts(service)).Name("get_all")
		router.Get("/host", r.mw.Authentication(scope.DhcpCanRead...), handler.GetStaticHost(service)).Name("get")
		router.Post("/host", r.mw.Authentication(scope.DhcpCanAdd...), handler.AddStaticHost(service)).Name("add")
		router.Put("/host", r.mw.Authentication(scope.DhcpCanChange...), handler.UpdateStaticHost(service)).Name("update")
		router.Delete("/host", r.mw.Authentication(scope.DhcpCanChange...), handler.RemoveStaticHost(service)).Name("remove")
	}, "static.hosts.")
}

func (r Router) Metrics(cfg monitor.Config) {
	r.root.Get("/metrics", monitor.New(cfg))
}

func (r Router) SwaggerUI(openApiSpecFile string) {
	fiberswagger.Router(r.root, fiberswagger.Config{
		BasePath: "/openapi",
		FilePath: openApiSpecFile,
	})
}
