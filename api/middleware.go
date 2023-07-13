package api

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gringolito/pi-hole-manager/api/middleware/fiberslog"
	"github.com/gringolito/pi-hole-manager/config"
	"golang.org/x/exp/slog"
)

type Middleware interface {
	Authentication(roles ...string) fiber.Handler
	Logger() fiber.Handler
	Recovery() fiber.Handler
	RequestId() fiber.Handler
}

func NewMiddleware(logger *slog.Logger, cfg *config.Config) (Middleware, error) {
	jwtConfig, err := setupJwtConfig(cfg)
	if err != nil {
		return nil, err
	}

	return middleware{
		logger: fiberslog.New(fiberslog.Config{
			Logger: logger,
			Fields: []string{"latency", "status", "method", "path", "requestId", "ip", "port", "pid"},
		}),
		recovery: recover.New(recover.Config{
			EnableStackTrace: true,
		}),
		requestId: requestid.New(),
		jwtConfig: jwtConfig,
	}, nil
}

type middleware struct {
	logger    fiber.Handler
	recovery  fiber.Handler
	requestId fiber.Handler
	jwtConfig *jwtware.Config
}

var voidMiddleware = func(c *fiber.Ctx) error {
	return c.Next()
}

func (m middleware) Authentication(roles ...string) fiber.Handler {
	if m.jwtConfig == nil {
		return voidMiddleware
	}

	contextKey := "user"
	if m.jwtConfig.ContextKey != "" {
		contextKey = m.jwtConfig.ContextKey
	}

	if len(roles) > 0 {
		m.jwtConfig.SuccessHandler = authorizationHandler(contextKey, roles)
	}

	return jwtware.New(*m.jwtConfig)
}

func (m middleware) Logger() fiber.Handler {
	return m.logger
}

func (m middleware) Recovery() fiber.Handler {
	return m.recovery
}

func (m middleware) RequestId() fiber.Handler {
	return m.requestId
}
