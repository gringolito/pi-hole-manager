package model

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strings"
)

type StaticDhcpHost struct {
	MacAddress net.HardwareAddr
	IPAddress  net.IP
	HostName   string
	Interface  string
}

var ErrInvalidDHCPHost = errors.New("invalid DHCP host entry")

func (h *StaticDhcpHost) FromConfig(config string) error {
	tokens := strings.Split(config, ",")
	if len(tokens) != 3 {
		return ErrInvalidDHCPHost
	}

	var mac string
	_, err := fmt.Sscanf(tokens[0], "dhcp-host=%s", &mac)
	if err != nil {
		return err
	}

	h.MacAddress, err = net.ParseMAC(mac)
	h.IPAddress = net.ParseIP(tokens[1])
	h.HostName = tokens[2]

	return err
}

func (h *StaticDhcpHost) ToConfig() string {
	return fmt.Sprintf("dhcp-host=%s,%s,%s", h.MacAddress.String(), h.IPAddress.String(), h.HostName)
}

func (h *StaticDhcpHost) Equal(other StaticDhcpHost) bool {
	return bytes.Equal(h.MacAddress, other.MacAddress) && bytes.Equal(h.IPAddress, other.IPAddress) && h.HostName == other.HostName
}
