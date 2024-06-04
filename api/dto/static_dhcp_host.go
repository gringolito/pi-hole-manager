package dto

import (
	"net"

	"github.com/gringolito/dnsmasq-manager/pkg/model"
)

type StaticDhcpHost struct {
	MacAddress string `validate:"required,mac"`
	IPAddress  string `validate:"required,ipv4"`
	HostName   string `validate:"required,hostname"`
	Interface  string `validate:"optional,printascii"`
}

func NewStaticDhcpHost(host *model.StaticDhcpHost) *StaticDhcpHost {
	return &StaticDhcpHost{
		MacAddress: host.MacAddress.String(),
		IPAddress:  host.IPAddress.String(),
		HostName:   host.HostName,
		Interface:  host.Interface,
	}
}

func (h *StaticDhcpHost) ToModel() *model.StaticDhcpHost {
	mac, _ := net.ParseMAC(h.MacAddress)

	return &model.StaticDhcpHost{
		MacAddress: mac,
		IPAddress:  net.ParseIP(h.IPAddress),
		HostName:   h.HostName,
		Interface:  h.Interface,
	}
}
