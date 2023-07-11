package api

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gringolito/pi-hole-manager/api/middleware/fiberslog"
	"github.com/gringolito/pi-hole-manager/config"
	"golang.org/x/exp/slog"
)

var voidMiddleware = func(c *fiber.Ctx) error {
	return c.Next()
}

type Middleware interface {
	Authentication() fiber.Handler
	Logger() fiber.Handler
	Recovery() fiber.Handler
	RequestId() fiber.Handler
}

func NewMiddleware(logger *slog.Logger, cfg *config.Config) (Middleware, error) {
	auth, err := setupAuthHandler(logger, cfg)
	if err != nil {
		return middleware{}, err
	}

	return middleware{
		authentication: auth,
		logger: fiberslog.New(fiberslog.Config{
			Logger: logger,
			Fields: []string{"latency", "status", "method", "path", "requestId", "ip", "port", "pid"},
		}),
		recovery: recover.New(recover.Config{
			EnableStackTrace: true,
		}),
		requestId: requestid.New(),
	}, nil
}

func setupAuthHandler(logger *slog.Logger, cfg *config.Config) (fiber.Handler, error) {
	if cfg.Auth.Method == config.NoAuth {
		return voidMiddleware, nil
	}

	signingMethod, pemEncoded, err := getAuthSigningMethod(cfg.Auth.Method)
	if err != nil {
		logger.Debug(
			"invalid authentication signing method, please check your configuration",
			slog.Group("config",
				slog.Group("auth",
					slog.String("method", cfg.Auth.Method))),
		)
		return nil, err
	}

	signingKey, err := getAuthSigningKey(cfg.Auth.Key, pemEncoded)
	if err != nil {
		logger.Debug(
			"invalid or malformed authentication signing key, please check your configuration",
			slog.Group("config",
				slog.Group("auth",
					slog.String("method", cfg.Auth.Method))),
		)
		return nil, err
	}

	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: signingMethod,
			Key:    signingKey,
		},
	}), nil
}

func getAuthSigningMethod(authMethod string) (string, bool, error) {
	type signMethod struct {
		method     string
		pemEncoded bool
	}

	signingMethods := map[string]*signMethod{
		config.AuthES256: {jwtware.ES256, true},
		config.AuthES384: {jwtware.ES384, true},
		config.AuthES512: {jwtware.ES512, true},
		config.AuthHS256: {jwtware.HS256, false},
		config.AuthHS384: {jwtware.HS384, false},
		config.AuthHS512: {jwtware.HS512, false},
		config.AuthRS256: {jwtware.RS256, true},
		config.AuthRS384: {jwtware.RS384, true},
		config.AuthRS512: {jwtware.RS512, true},
	}

	signing := signingMethods[authMethod]
	if signing == nil {
		return "", false, fmt.Errorf("invalid auth signing method: %s", authMethod)
	}

	return signing.method, signing.pemEncoded, nil
}

func getAuthSigningKey(authKey string, pemEncoded bool) (interface{}, error) {
	if !pemEncoded {
		// HMAC are just plain text
		return []byte(authKey), nil
	}

	key, err := os.ReadFile(authKey)
	if err != nil {
		key = []byte(authKey)
	}

	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("invalid or malformed auth signing key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

type middleware struct {
	authentication fiber.Handler
	logger         fiber.Handler
	recovery       fiber.Handler
	requestId      fiber.Handler
}

func (m middleware) Authentication() fiber.Handler {
	return m.authentication
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
