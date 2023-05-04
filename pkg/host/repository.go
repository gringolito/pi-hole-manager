package host

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/gringolito/pi-hole-manager/pkg/model"
)

type Repository interface {
	Delete(host *model.StaticDhcpHost) (*model.StaticDhcpHost, error)
	DeleteByMac(macAddress string) (*model.StaticDhcpHost, error)
	DeleteByIP(ipAddress string) (*model.StaticDhcpHost, error)
	Find(host *model.StaticDhcpHost) (*model.StaticDhcpHost, error)
	FindAll() (*[]model.StaticDhcpHost, error)
	FindByMac(macAddress string) (*model.StaticDhcpHost, error)
	FindByIP(ipAddress string) (*model.StaticDhcpHost, error)
	Save(host *model.StaticDhcpHost) error
}

type repository struct {
	staticHostsFilePath string
}

func NewRepository(staticHostsFilePath string) Repository {
	return &repository{
		staticHostsFilePath: staticHostsFilePath,
	}
}

func (r *repository) FindAll() (*[]model.StaticDhcpHost, error) {
	return r.load()
}

func (r *repository) Find(host *model.StaticDhcpHost) (*model.StaticDhcpHost, error) {
	return r.find(func(other model.StaticDhcpHost) bool {
		return *host == other
	})
}

func (r *repository) FindByMac(macAddress string) (*model.StaticDhcpHost, error) {
	return r.find(func(other model.StaticDhcpHost) bool {
		return macAddress == other.MacAddress
	})
}

func (r *repository) FindByIP(ipAddress string) (*model.StaticDhcpHost, error) {
	return r.find(func(other model.StaticDhcpHost) bool {
		return ipAddress == other.IPAddress
	})
}

func (r *repository) Save(host *model.StaticDhcpHost) error {
	hosts, err := r.load()
	if err != nil {
		return err
	}

	*hosts = append(*hosts, *host)
	return r.save(hosts)
}

func (r *repository) Delete(host *model.StaticDhcpHost) (*model.StaticDhcpHost, error) {
	return r.delete(func(other model.StaticDhcpHost) bool {
		return *host == other
	})
}

func (r *repository) DeleteByMac(macAddress string) (*model.StaticDhcpHost, error) {
	return r.delete(func(other model.StaticDhcpHost) bool {
		return macAddress == other.MacAddress
	})
}

func (r *repository) DeleteByIP(ipAddress string) (*model.StaticDhcpHost, error) {
	return r.delete(func(other model.StaticDhcpHost) bool {
		return ipAddress == other.IPAddress
	})
}

func (r *repository) load() (*[]model.StaticDhcpHost, error) {
	file, err := os.Open(r.staticHostsFilePath)
	if err != nil {
		log.Printf("Error reading static hosts file (%s): %s", r.staticHostsFilePath, err.Error())
		return nil, err
	}
	defer file.Close()

	return r.parse(file)
}

func (r *repository) parse(file *os.File) (*[]model.StaticDhcpHost, error) {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	hosts := []model.StaticDhcpHost{}
	for scanner.Scan() {
		if !strings.HasPrefix(scanner.Text(), "dhcp-host=") {
			log.Println("Skipping line:", scanner.Text())
			continue
		}
		log.Println("Parsing line:", scanner.Text())

		host := model.StaticDhcpHost{}
		err := host.FromConfig(scanner.Text())
		if err != nil {
			log.Printf("Failed to parse static DHCP host entry (%s): %s", scanner.Text(), err.Error())
			return nil, err
		}

		hosts = append(hosts, host)
	}

	return &hosts, nil
}

func (r *repository) save(hosts *[]model.StaticDhcpHost) error {
	config := make([]string, 0, len(*hosts))
	for _, host := range *hosts {
		config = append(config, host.ToConfig())
	}

	return os.WriteFile(r.staticHostsFilePath, []byte(strings.Join(config, "\n")), os.FileMode(0644))
}

func (r *repository) delete(filter func(host model.StaticDhcpHost) bool) (*model.StaticDhcpHost, error) {
	hosts, err := r.load()
	if err != nil {
		return nil, err
	}

	h := *hosts
	for i, host := range h {
		if !filter(host) {
			continue
		}

		*hosts = append(h[:i], h[i+1:]...)
		err := r.save(hosts)
		return &host, err
	}

	return nil, nil
}

func (r *repository) find(filter func(host model.StaticDhcpHost) bool) (*model.StaticDhcpHost, error) {
	hosts, err := r.load()
	if err != nil {
		return nil, err
	}

	for _, host := range *hosts {
		if filter(host) {
			return &host, nil
		}
	}

	return nil, nil
}
