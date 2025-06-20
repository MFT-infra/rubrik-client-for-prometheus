package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config holds the configuration values for the Rubrik exporter.
type Config struct {
	RubrikIP       string `mapstructure:"rubrik_ip"`
	Username       string `mapstructure:"username"`
	Password       string `mapstructure:"password"`
	ApiToken       string `mapstructure:"api_token"`
	ServiceID      string `mapstructure:"service_id"`
	ServiceSecret  string `mapstructure:"service_secret"`
	PrometheusPort string `mapstructure:"prometheus_port"`
}


// LoadConfig loads the configuration from a file, CMD line arguments or environment variables.
// Package config provides functionality to load and manage configuration settings for the Rubrik exporter.
// It supports loading from a YAML file, JSON file, environment variables, cmd arguments and provides default values.

func LoadConfig() (*Config, error) {
	v := viper.New()

	// Set config file name (without extension) and add search paths
	v.SetConfigName("config") // Will look for config.yaml or config.json
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/rubrik-exporter/")
	v.AddConfigPath("$HOME/.rubrik-exporter")

	// Support ENV variables
	v.AutomaticEnv()
	v.SetEnvPrefix("rubrik")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set default values
	v.SetDefault("prometheus_port", "8080")

	// Add the directory of the running executable
	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)
		v.AddConfigPath(execDir)
	}

	// Try loading config.yaml or config.json
	var cfg Config
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("No config file found or failed to read: %v (continuing with env/flags)\n", err)
	} else {
		if err := v.Unmarshal(&cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	return &cfg, nil
}
