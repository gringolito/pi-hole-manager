package fiberslog

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
)

type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// SkipBody defines a function to skip log  "body" field when returned true.
	//
	// Optional. Default: nil
	SkipBody func(c *fiber.Ctx) bool

	// SkipResBody defines a function to skip log  "resBody" field when returned true.
	//
	// Optional. Default: nil
	SkipResBody func(c *fiber.Ctx) bool

	// GetResBody defines a function to get ResBody.
	//  eg: when use compress middleware, resBody is unreadable. you can set GetResBody func to get readable resBody.
	//
	// Optional. Default: nil
	GetResBody func(c *fiber.Ctx) []byte

	// Skip logging for these uri
	//
	// Optional. Default: nil
	SkipURIs []string

	// Add custom slog logger.
	//
	// Optional. Default: slog.New(slog.NewJSONHandler())
	Logger *slog.Logger

	// Add fields what you want see.
	//
	// Optional. Default: {"latency", "status", "method", "url", "pid"}
	Fields []string
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	Logger: slog.Default(),
	Fields: []string{"latency", "status", "method", "url", "pid"},
}

// Helper function to set default values
func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	// Set default values
	if cfg.Logger == nil {
		cfg.Logger = ConfigDefault.Logger
	}

	if cfg.Fields == nil {
		cfg.Fields = ConfigDefault.Fields
	}

	return cfg
}
