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

// New returns a fiber.Handler (middleware) that renders OpenAPI specification using SwaggerUI.
// New creates a new middleware handler
func New(config ...Config) fiber.Handler {
	// Set default config
	cfg := configDefault(config...)

	if _, err := os.Stat(cfg.FilePath); os.IsNotExist(err) {
		panic(fmt.Errorf("%s file is not exist", cfg.FilePath))
	}

	specDoc, err := loads.Spec(cfg.FilePath)
	if err != nil {
		panic(err)
	}

	spec, err := json.Marshal(specDoc.Raw())
	if err != nil {
		panic(err)
	}

	return adaptor.HTTPMiddleware(func(next http.Handler) http.Handler {
		swaggerUiHandler := middleware.SwaggerUI(middleware.SwaggerUIOpts{
			Path:    cfg.BasePath,
			SpecURL: path.Join(cfg.BasePath, "swagger.json"),
		}, next)

		return middleware.Spec(cfg.BasePath, spec, swaggerUiHandler)
	})
}
