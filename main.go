package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type staticDhcpHost struct {
	MacAddress string `json:"mac_address"`
	IpAddress  string `json:"ip_address"`
	HostName   string `json:"hostname"`
}

var hosts = []staticDhcpHost{
	{MacAddress: "0A:50:E2:E9:02:22", IpAddress: "192.168.11.145", HostName: "docker02"},
	{MacAddress: "42:E0:F0:A2:01:93", IpAddress: "192.168.11.146", HostName: "docker03"},
	{MacAddress: "80:3F:5D:DB:7B:53", IpAddress: "192.168.11.4", HostName: "wavlink-03"},
	{MacAddress: "7C:C7:09:E2:D0:6B", IpAddress: "192.168.11.126", HostName: "chip-wifi"},
}

func main() {
	router := gin.Default()
	router.GET("/static/hosts", getStaticHosts)
	router.GET("/static/host/:mac_address", getStaticHostByMac)
	router.POST("/static/hosts", addStaticHost)

	router.Run("localhost:8080")
}

func getStaticHosts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, hosts)
}

func addStaticHost(c *gin.Context) {
	var newHost staticDhcpHost

	if err := c.BindJSON(&newHost); err != nil {
		return
	}

	hosts = append(hosts, newHost)
	c.IndentedJSON(http.StatusCreated, newHost)
}

func getStaticHostByMac(c *gin.Context) {
	macAddress := c.Param("mac_address")

	for _, h := range hosts {
		if h.MacAddress == macAddress {
			c.IndentedJSON(http.StatusOK, h)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "host not found"})
}
