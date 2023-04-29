package hosts

import (
	"fmt"
	"strings"
)

type staticDhcpHost struct {
	MacAddress string `json:"mac_address"`
	IPAddress  string `json:"ip_address"`
	HostName   string `json:"hostname"`
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
