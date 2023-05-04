package handler

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gringolito/pi-hole-manager/pkg/host"
	"github.com/gringolito/pi-hole-manager/pkg/model"
	"github.com/gringolito/pi-hole-manager/pkg/validation"
)

func getHostFromBody(c *fiber.Ctx) *model.StaticDhcpHost {
	host := new(model.StaticDhcpHost)
	if err := c.BodyParser(host); err != nil {
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		return nil
	}

	if errors := validation.Validate(host); errors != nil {
		c.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": errors})
		return nil
	}

	return host
}

func GetAllStaticHosts(service host.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		hosts, err := service.FetchAll()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(http.StatusOK).JSON(hosts)
	}
}

func GetStaticHost(service host.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		macAddress := c.Query("mac")
		if len(macAddress) > 0 {
			return getStaticHostByMac(service, c, macAddress)
		}

		ipAddress := c.Query("ip")
		if len(ipAddress) > 0 {
			return getStaticHostByIP(service, c, ipAddress)
		}

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "No query parameter specified"})
	}
}

func getStaticHostByMac(service host.Service, c *fiber.Ctx, macAddress string) error {
	host, err := service.FetchByMac(macAddress)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if host == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": fmt.Sprintf("MAC address '%s' not found", macAddress)})
	}

	return c.Status(http.StatusOK).JSON(host)
}

func getStaticHostByIP(service host.Service, c *fiber.Ctx, ipAddress string) error {
	host, err := service.FetchByIP(ipAddress)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if host == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": fmt.Sprintf("IP address '%s' not found", ipAddress)})
	}

	return c.Status(http.StatusOK).JSON(host)
}

func AddStaticHost(service host.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		host := getHostFromBody(c)
		if host == nil {
			return nil
		}

		if err := service.Insert(host); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(http.StatusCreated).JSON(host)
	}
}

func UpdateStaticHost(service host.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		host := getHostFromBody(c)
		if host == nil {
			return nil
		}

		if err := service.Update(host); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(http.StatusCreated).JSON(host)
	}
}

func RemoveStaticHost(service host.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		macAddress := c.Query("mac")
		if len(macAddress) > 0 {
			return removeStaticHostByMac(service, c, macAddress)
		}

		ipAddress := c.Query("ip")
		if len(ipAddress) > 0 {
			return removeStaticHostByIP(service, c, ipAddress)
		}

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "No query parameter specified"})
	}
}

func removeStaticHostByMac(service host.Service, c *fiber.Ctx, macAddress string) error {
	host, err := service.RemoveByMac(macAddress)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(host)
}

func removeStaticHostByIP(service host.Service, c *fiber.Ctx, ipAddress string) error {
	host, err := service.RemoveByIP(ipAddress)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(host)
}
