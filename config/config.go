package config

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	LogLevelDebug   = "debug"
	LogLevelInfo    = "info"
	LogLevelWarning = "warning"
	LogLevelError   = "error"
)

const (
	LogFormatJSON      = "json"
	LogFormatPlainText = "text"
)

// Other default constants
const (
	DefaultDhcpStaticHostFile = "/etc/dnsmasq.d/04-pihole-static-dhcp.conf"
	DefaultServerHttpPort     = 6904
)

type Config struct {
	Auth struct {
		Username string
		Password string
	}
	Host struct {
		Static struct {
			File string
		}
	}
	Server struct {
		Port int
	}
	Log struct {
		Level  string
		File   string
		Format string
		Source bool
	}
}

func newDefaultConfig() *Config {
	def := Config{}
	def.Host.Static.File = "/etc/dnsmasq.d/04-pihole-static-dhcp.conf"
	def.Server.Port = 6904
	def.Log.Level = LogLevelInfo
	def.Log.Format = LogFormatJSON

	return &def
}

func Init(configName string) (*Config, error) {
	viper.AddConfigPath("/etc/pi-hole-monitor/")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName(configName)
	viper.SetEnvPrefix("PHM") // PHM stands for Pi-Hole Manager
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error, parse environments and load defaults
		} else {
			return nil, err
		}
	}

	config := newDefaultConfig()
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	return config, err
}
