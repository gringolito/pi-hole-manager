package hosts

import "fmt"

type hostService struct {
	repository hostRepository
}

func NewService(staticHostsFilePath string) hostService {
	return hostService{repository: NewRepository(staticHostsFilePath)}
}

func (s *hostService) AddStaticHost(host staticDhcpHost) error {
	sameMacHost, err := s.repository.FindByMac(host.MacAddress)
	if err != nil {
		return err
	}
	if sameMacHost != nil {
		return fmt.Errorf("Duplicated MAC address")
	}

	sameIPHost, err := s.repository.FindByIP(host.IPAddress)
	if err != nil {
		return err
	}
	if sameIPHost != nil {
		return fmt.Errorf("Duplicated IP address")
	}

	return s.repository.Insert(host)
}

func (s *hostService) ForceAddStaticHost(host staticDhcpHost) error {
	_, err := s.repository.RemoveByMac(host.MacAddress)
	if err != nil {
		return err
	}

	_, err = s.repository.RemoveByIP(host.IPAddress)
	if err != nil {
		return err
	}

	return s.repository.Insert(host)
}

func (s *hostService) GetStaticHosts() ([]staticDhcpHost, error) {
	return s.repository.Load()
}

func (s *hostService) GetStaticHostByMac(macAddress string) (*staticDhcpHost, error) {
	return s.repository.FindByMac(macAddress)
}

func (s *hostService) GetStaticHostByIP(ipAddress string) (*staticDhcpHost, error) {
	return s.repository.FindByIP(ipAddress)
}

func (s *hostService) DeleteStaticHostByMac(macAddress string) (*staticDhcpHost, error) {
	return s.repository.RemoveByMac(macAddress)
}

func (s *hostService) DeleteStaticHostByIP(ipAddress string) (*staticDhcpHost, error) {
	return s.repository.RemoveByIP(ipAddress)
}
