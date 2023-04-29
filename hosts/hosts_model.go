package hosts

import (
	"fmt"
	"strings"
)

type staticDhcpHost struct {
	MacAddress string `json:"mac_address" binding:"required,mac"`
	IPAddress  string `json:"ip_address" binding:"required,ipv4"`
	HostName   string `json:"hostname" binding:"required,hostname"`
}

func (h *staticDhcpHost) FromConfig(config string) error {
	tokens := strings.Split(config, ",")
	if len(tokens) != 3 {
		return fmt.Errorf("Invalid DHCP host entry")
	}

	_, err := fmt.Sscanf(tokens[0], "dhcp-host=%s", &h.MacAddress)
	if err != nil {
		return err
	}
	h.IPAddress = tokens[1]
	h.HostName = tokens[2]

	return nil
}

func (h *staticDhcpHost) ToConfig() string {
	return fmt.Sprintf("dhcp-host=%s,%s,%s", h.MacAddress, h.IPAddress, h.HostName)
}
