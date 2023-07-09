package fiberswagger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func loadSpec(cfg Config) *loads.Document {
	if _, err := os.Stat(cfg.FilePath); os.IsNotExist(err) {
		panic(fmt.Errorf("%s file is not exist", cfg.FilePath))
	}

	spec, err := loads.Spec(cfg.FilePath)
	if err != nil {
		panic(err)
	}

	return spec
}

// Middleware returns a fiber.Handler (middleware) that renders OpenAPI specification using SwaggerUI.
// Middleware creates a new middleware handler
func Middleware(config ...Config) fiber.Handler {
	// Set default config
	cfg := configDefault(config...)

	spec := loadSpec(cfg)

	specJson, err := json.Marshal(spec.Raw())
	if err != nil {
		panic(err)
	}

	return adaptor.HTTPMiddleware(func(next http.Handler) http.Handler {
		swaggerUiHandler := middleware.SwaggerUI(middleware.SwaggerUIOpts{
			Path:    cfg.BasePath,
			SpecURL: path.Join(cfg.BasePath, "swagger.json"),
		}, next)

		return middleware.Spec(cfg.BasePath, specJson, swaggerUiHandler)
	})
}

// Router creates routes with handlers to renders OpenAPI specification using SwaggerUI.
func Router(router fiber.Router, config ...Config) {
	// Set default config
	cfg := configDefault(config...)

	spec := loadSpec(cfg)

	router.Route(cfg.BasePath, func(router fiber.Router) {
		router.Get("/", handleSwaggerUi(cfg)).Name("ui")
		router.Get("/swagger.json", handleSwaggerJson(spec.Raw())).Name("spec")
	}, "swagger.")
}

func handleSwaggerUi(cfg Config) fiber.Handler {
	return adaptor.HTTPHandler(middleware.SwaggerUI(middleware.SwaggerUIOpts{
		Path:    cfg.BasePath,
		SpecURL: path.Join(cfg.BasePath, "swagger.json"),
	}, nil))
}

func handleSwaggerJson(swagger json.RawMessage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Status(http.StatusOK).JSON(swagger)
		return nil
	}
}
