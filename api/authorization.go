package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"
)

func authorizationHandler(jwtContextKey string, roles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Locals(jwtContextKey).(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		name := claims["name"].(string)

		groups, err := parseClaimGroups(claims)
		if err != nil {
			slog.Debug("authorization denied: malformed JWT claim 'groups'", slog.String("user", name))
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error":   "Forbidden",
				"message": "You do not have permission to access this resource.",
				"details": "The user provided a malformed JWT.",
			})
		}

		authorized := false
		for _, role := range roles {
			if slices.Contains[string](groups, role) {
				authorized = true
				slog.Debug("authorization granted", slog.String("user", name), slog.String("role", role))
				break
			}
		}

		if !authorized {
			slog.Debug("authorization denied: required role not found", slog.String("user", name))
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error":   "Forbidden",
				"message": "You do not have permission to access this resource.",
				"details": "The user does not have the required role to access this resource.",
			})
		}

		return c.Next()
	}
}

func parseClaimGroups(claims jwt.MapClaims) ([]string, error) {
	var groups []string

	switch v := claims["groups"].(type) {
	case string:
		groups = append(groups, v)
	case []string:
		groups = v
	case []interface{}:
		for _, a := range v {
			group, ok := a.(string)
			if !ok {
				return nil, jwt.ErrInvalidType
			}
			groups = append(groups, group)
		}
	}

	return groups, nil
}
