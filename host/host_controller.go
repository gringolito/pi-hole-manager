package host

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gringolito/pi-hole-manager/validation"
)

type Controller struct {
	service hostService
}

func NewController(staticHostsFilePath string) Controller {
	return Controller{service: NewService(staticHostsFilePath)}
}

func (h *Controller) SetupRoutes(router fiber.Router) {
	static := router.Group("/static")
	static.Get("/hosts", h.getStaticHosts).Name("Get All Static Hosts")
	static.Get("/host", h.getStaticHost).Name("Find Static Host")
	static.Post("/host", h.addStaticHost).Name("Create Static Host")
	static.Put("/host", h.addOrReplaceStaticHost).Name("Create or Replace Static Host")
	static.Delete("/host", h.deleteStaticHost).Name("Remove Static Host")
}

func (h *Controller) getHostFromBody(c *fiber.Ctx) *staticDhcpHost {
	host := new(staticDhcpHost)
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

func (h *Controller) getStaticHosts(c *fiber.Ctx) error {
	hosts, err := h.service.GetStaticHosts()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(hosts)
}

func (h *Controller) addStaticHost(c *fiber.Ctx) error {
	host := h.getHostFromBody(c)
	if host == nil {
		return nil
	}

	if err := h.service.AddStaticHost(host); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(host)
}

func (h *Controller) addOrReplaceStaticHost(c *fiber.Ctx) error {
	host := h.getHostFromBody(c)
	if host == nil {
		return nil
	}

	if err := h.service.ForceAddStaticHost(host); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(host)
}
func (h *Controller) getStaticHost(c *fiber.Ctx) error {
	macAddress := c.Query("mac")
	if len(macAddress) > 0 {
		h.getStaticHostByMac(c, macAddress)
	}

	ipAddress := c.Query("ip")
	if len(ipAddress) > 0 {
		h.getStaticHostByIP(c, ipAddress)
	}

	return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "No query parameter specified"})
}

func (h *Controller) getStaticHostByMac(c *fiber.Ctx, macAddress string) error {
	host, err := h.service.GetStaticHostByMac(macAddress)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if host == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": fmt.Sprintf("MAC address '%s' not found", macAddress)})
	}

	return c.Status(http.StatusOK).JSON(host)
}

func (h *Controller) getStaticHostByIP(c *fiber.Ctx, ipAddress string) error {
	host, err := h.service.GetStaticHostByIP(ipAddress)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if host == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": fmt.Sprintf("IP address '%s' not found", ipAddress)})
	}

	return c.Status(http.StatusOK).JSON(host)
}

func (h *Controller) deleteStaticHost(c *fiber.Ctx) error {
	macAddress := c.Query("mac")
	if len(macAddress) > 0 {
		return h.deleteStaticHostByMac(c, macAddress)
	}

	ipAddress := c.Query("ip")
	if len(ipAddress) > 0 {
		return h.deleteStaticHostByIP(c, ipAddress)
	}

	return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "No query parameter specified"})
}

func (h *Controller) deleteStaticHostByMac(c *fiber.Ctx, macAddress string) error {
	host, err := h.service.DeleteStaticHostByMac(macAddress)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(host)
}

func (h *Controller) deleteStaticHostByIP(c *fiber.Ctx, ipAddress string) error {
	host, err := h.service.DeleteStaticHostByIP(ipAddress)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(host)
}
