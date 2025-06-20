package main

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

// Config holds the configuration values for the Rubrik exporter.
type Config struct {
	RubrikIP       string `json:"rubrik_ip" yaml:"rubrik_ip"`
	Username       string `json:"username" yaml:"username"`
	Password       string `json:"password" yaml:"password"`
	ApiToken       string `json:"api_token" yaml:"api_token"`
	ServiceID      string `json:"service_id" yaml:"service_id"`
	ServiceSecret  string `json:"service_secret" yaml:"service_secret"`
	PrometheusPort string `json:"prometheus_port" yaml:"prometheus_port"`
}

// LoadConfig loads the configuration from a file, CMD line arguments or environment variables.
// Package config provides functionality to load and manage configuration settings for the Rubrik exporter.
// It supports loading from a YAML file, JSON file, environment variables, cmd arguments and provides default values.

func LoadConfig() (*Config, error) {
	v := viper.New()

	// Allow config from file
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/rubrik-exporter/")
	v.AddConfigPath("$HOME/.rubrik-exporter")
	v.SetConfigType("yaml")

	// Check for both JSON and YAML
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("No config file found: %v (continuing)\n", err)
	}

	// Support ENV variables like RUBRIK_IP, etc.
	v.AutomaticEnv()
	v.SetEnvPrefix("rubrik") // prefix like RUBRIK_USERNAME
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Default fallback port
	v.SetDefault("prometheus_port", "8080")

	err := v.ReadInConfig()
	if err != nil {
		fmt.Printf("No config file found: %v (continuing)\n", err)
	}

	var cfg Config
	err = v.Unmarshal(&cfg)

	var (
		flagRubrikIP       = flag.String("rubrik-ip", "", "Rubrik IP or FQDN")
		flagUsername       = flag.String("username", "", "Rubrik username")
		flagPassword       = flag.String("password", "", "Rubrik password")
		flagApiToken       = flag.String("api-token", "", "Rubrik API token")
		flagServiceID      = flag.String("service-id", "", "Service account ID")
		flagServiceSecret  = flag.String("service-secret", "", "Service account secret")
		flagPrometheusPort = flag.String("prometheus-port", "", "Prometheus HTTP port")
	)
	flag.Parse()

	// Override from flags if set
	if *flagRubrikIP != "" {
		cfg.RubrikIP = *flagRubrikIP
	}
	if *flagUsername != "" {
		cfg.Username = *flagUsername
	}
	if *flagPassword != "" {
		cfg.Password = *flagPassword
	}
	if *flagApiToken != "" {
		cfg.ApiToken = *flagApiToken
	}
	if *flagServiceID != "" {
		cfg.ServiceID = *flagServiceID
	}
	if *flagServiceSecret != "" {
		cfg.ServiceSecret = *flagServiceSecret
	}
	if *flagPrometheusPort != "" {
		cfg.PrometheusPort = *flagPrometheusPort
	}

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
