package hosts

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service hostService
}

func NewController(staticHostsFilePath string) Controller {
	return Controller{service: NewService(staticHostsFilePath)}
}

func (h *Controller) SetupRoutes(router *gin.Engine) {
	router.GET("/static/hosts", h.getStaticHosts)
	router.GET("/static/host", h.getStaticHost)
	router.POST("/static/host", h.addStaticHost)
	router.PUT("/static/host", h.addOrReplaceStaticHost)
	router.DELETE("/static/host", h.deleteStaticHost)
}

func (h *Controller) getStaticHosts(c *gin.Context) {
	hosts, err := h.service.GetStaticHosts()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, hosts)
}

func (h *Controller) addStaticHost(c *gin.Context) {
	var newHost staticDhcpHost
	if err := c.ShouldBindJSON(&newHost); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := h.service.AddStaticHost(newHost); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, newHost)
}

func (h *Controller) addOrReplaceStaticHost(c *gin.Context) {
	var newHost staticDhcpHost
	if err := c.ShouldBindJSON(&newHost); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := h.service.ForceAddStaticHost(newHost); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, newHost)
}
func (h *Controller) getStaticHost(c *gin.Context) {
	macAddress, ok := c.GetQuery("mac")
	if ok {
		h.getStaticHostByMac(c, macAddress)
		return
	}

	ipAddress, ok := c.GetQuery("ip")
	if ok {
		h.getStaticHostByIP(c, ipAddress)
		return
	}

	c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "No query parameter specified"})
}

func (h *Controller) getStaticHostByMac(c *gin.Context, macAddress string) {
	host, err := h.service.GetStaticHostByMac(macAddress)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if host == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("MAC address '%s' not found", macAddress)})
		return
	}

	c.IndentedJSON(http.StatusOK, host)
}

func (h *Controller) getStaticHostByIP(c *gin.Context, ipAddress string) {
	host, err := h.service.GetStaticHostByIP(ipAddress)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if host == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("IP address '%s' not found", ipAddress)})
		return
	}

	c.IndentedJSON(http.StatusOK, host)
}

func (h *Controller) deleteStaticHost(c *gin.Context) {
	macAddress, ok := c.GetQuery("mac")
	if ok {
		h.deleteStaticHostByMac(c, macAddress)
		return
	}

	ipAddress, ok := c.GetQuery("ip")
	if ok {
		h.deleteStaticHostByIP(c, ipAddress)
		return
	}

	c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "No query parameter specified"})
}

func (h *Controller) deleteStaticHostByMac(c *gin.Context, macAddress string) {
	host, err := h.service.DeleteStaticHostByMac(macAddress)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, host)
}

func (h *Controller) deleteStaticHostByIP(c *gin.Context, ipAddress string) {
	host, err := h.service.DeleteStaticHostByIP(ipAddress)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, host)
}
