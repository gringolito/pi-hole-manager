package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gringolito/pi-hole-manager/hosts"
)

const STATIC_HOSTS_FILE_PATH string = "./04-pihole-static-dhcp.conf"

func main() {
	router := gin.Default()

	h := hosts.NewController(STATIC_HOSTS_FILE_PATH)
	h.SetupRoutes(router)

	router.Run("localhost:8080")
}
