package config

import (
	"strings"

	"github.com/spf13/viper"
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
		Level string
	}
}

func Init(configName string) (*Config, error) {
	// Read the defaults first
	viper.AddConfigPath("/etc/pi-hole-monitor/")
	viper.AddConfigPath("config/")
	viper.SetConfigName("default")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	// Than merge with the specifics
	viper.SetConfigName(configName)
	viper.SetEnvPrefix("PHM") // PHM stands for Pi-Hole Manager
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	err = viper.MergeInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error and parse environments
		} else {
			return nil, err
		}
	}

	config := new(Config)
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	return config, err
}
