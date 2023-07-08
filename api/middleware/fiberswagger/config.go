package fiberswagger

type Config struct {
	// Base path for the SwaggerUI.
	//
	// Optional. Default: "/docs"
	BasePath string

	// OpenAPI specification file path to be rendered.
	//
	// Optional. Default: "./openapi.yaml"
	FilePath string
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	BasePath: "/docs",
	FilePath: "./openapi.yaml",
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
	if cfg.BasePath == "" {
		cfg.BasePath = ConfigDefault.BasePath
	}

	if cfg.FilePath == "" {
		cfg.FilePath = ConfigDefault.FilePath
	}

	return cfg
}
