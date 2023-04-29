package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gringolito/pi-hole-manager/hosts"
)

const STATIC_HOSTS_FILE_PATH string = "./04-pihole-static-dhcp.conf"

func main() {
	router := gin.Default()

	hosts.Init(router, STATIC_HOSTS_FILE_PATH)

	router.Run("localhost:8080")
}
