package hosts

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service hostService
}

func Init(router *gin.Engine, staticHostsFilePath string) Controller {
	c := Controller{service: NewService(staticHostsFilePath)}
	c.setupRoutes(router)
	return c
}

func (h *Controller) setupRoutes(router *gin.Engine) {
	router.GET("/static/hosts", h.getStaticHosts)
	router.GET("/static/host/:mac_address", h.getStaticHostByMac)
	router.POST("/static/hosts", h.addStaticHost)
}

func (h *Controller) getStaticHosts(c *gin.Context) {
	hosts, err := h.service.getStaticHosts()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, hosts)
}

func (h *Controller) addStaticHost(c *gin.Context) {
	var newHost staticDhcpHost

	if err := c.BindJSON(&newHost); err != nil {
		return
	}

	// hosts = append(hosts, newHost)
	c.IndentedJSON(http.StatusCreated, newHost)
}

func (h *Controller) getStaticHostByMac(c *gin.Context) {
	macAddress := c.Param("mac_address")

	hosts, _ := h.service.getStaticHosts()
	for _, host := range hosts {
		if host.MacAddress == macAddress {
			c.IndentedJSON(http.StatusOK, host)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "host not found"})
}
