package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Auth.Method constants
const (
	NoAuth    = "none"
	AuthES256 = "ecdsa-256"
	AuthES384 = "ecdsa-384"
	AuthES512 = "ecdsa-512"
	AuthHS256 = "hmac-256"
	AuthHS384 = "hmac-384"
	AuthHS512 = "hmac-512"
	AuthRS256 = "rsa-256"
	AuthRS384 = "rsa-384"
	AuthRS512 = "rsa-512"
)

// Log.Level constants
const (
	LogLevelDebug   = "debug"
	LogLevelInfo    = "info"
	LogLevelWarning = "warning"
	LogLevelError   = "error"
)

// Log.Format constants
const (
	LogFormatJSON      = "json"
	LogFormatPlainText = "text"
)

// Other default constants
const (
	DefaultDhcpStaticHostFile = "/etc/dnsmasq.d/04-dhcp-static-leases.conf"
	DefaultServerHttpPort     = 6904
)

type Config struct {
	Auth struct {
		Method string
		Key    string
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
	def.Auth.Method = NoAuth
	def.Host.Static.File = DefaultDhcpStaticHostFile
	def.Server.Port = DefaultServerHttpPort
	def.Log.Level = LogLevelInfo
	def.Log.Format = LogFormatJSON

	return &def
}

func Init(configName string) (*Config, error) {
	viper.AddConfigPath("/etc/dnsmasq-manager/")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName(configName)
	viper.SetEnvPrefix("DMM") // DMM stands for (d)ns(m)asq (M)anager
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
