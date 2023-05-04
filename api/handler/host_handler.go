package handler

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gringolito/pi-hole-manager/api/dto"
	"github.com/gringolito/pi-hole-manager/api/validation"
	"github.com/gringolito/pi-hole-manager/pkg/host"
	"github.com/gringolito/pi-hole-manager/pkg/model"
)

func getHostFromBody(c *fiber.Ctx) *model.StaticDhcpHost {
	host := new(dto.StaticDhcpHost)
	if err := c.BodyParser(host); err != nil {
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		return nil
	}

	if errors := validation.Validate(host); errors != nil {
		c.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": errors})
		return nil
	}

	return host.ToModel()
}

func toStaticDhcpHostsDto(hosts *[]model.StaticDhcpHost) *[]dto.StaticDhcpHost {
	response := make([]dto.StaticDhcpHost, 0, len(*hosts))
	for _, h := range *hosts {
		response = append(response, *dto.NewStaticDhcpHost(&h))
	}

	return &response
}

func GetAllStaticHosts(service host.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		hosts, err := service.FetchAll()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(http.StatusOK).JSON(toStaticDhcpHostsDto(hosts))
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
	mac, err := net.ParseMAC(macAddress)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	host, err := service.FetchByMac(mac)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if host == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": fmt.Sprintf("MAC address '%s' not found", macAddress)})
	}

	return c.Status(http.StatusOK).JSON(dto.NewStaticDhcpHost(host))
}

func getStaticHostByIP(service host.Service, c *fiber.Ctx, ipAddress string) error {
	host, err := service.FetchByIP(net.ParseIP(ipAddress))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if host == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": fmt.Sprintf("IP address '%s' not found", ipAddress)})
	}

	return c.Status(http.StatusOK).JSON(dto.NewStaticDhcpHost(host))
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

		return c.Status(http.StatusCreated).JSON(dto.NewStaticDhcpHost(host))
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

		return c.Status(http.StatusCreated).JSON(dto.NewStaticDhcpHost(host))
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
	mac, err := net.ParseMAC(macAddress)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	host, err := service.RemoveByMac(mac)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if host == nil {
		return c.SendStatus(http.StatusNoContent)
	}

	return c.Status(http.StatusOK).JSON(dto.NewStaticDhcpHost(host))
}

func removeStaticHostByIP(service host.Service, c *fiber.Ctx, ipAddress string) error {
	host, err := service.RemoveByIP(net.ParseIP(ipAddress))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if host == nil {
		return c.SendStatus(http.StatusNoContent)
	}

	return c.Status(http.StatusOK).JSON(dto.NewStaticDhcpHost(host))
}
