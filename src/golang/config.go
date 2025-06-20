type Config struct {
    RubrikIP       string `json:"rubrik_ip" yaml:"rubrik_ip"`
    Username       string `json:"username" yaml:"username"`
    Password       string `json:"password" yaml:"password"`
    ApiToken       string `json:"api_token" yaml:"api_token"`
    ServiceID      string `json:"service_id" yaml:"service_id"`
    ServiceSecret  string `json:"service_secret" yaml:"service_secret"`
    PrometheusPort string `json:"prometheus_port" yaml:"prometheus_port"`
}
package config

import (
	"fmt"
	"strings"
	"github.com/spf13/viper"
)
// config.go
// Package config provides functionality to load and manage configuration settings for the Rubrik exporter.
// It supports loading from a YAML file, JSON file, environment variables, and provides default values.


func LoadConfig() (*Config, error) {
    v := viper.New()

    // Allow config from file
    v.SetConfigName("config") // no extension
    v.SetConfigType("yaml")
    v.AddConfigPath(".")
    v.AddConfigPath("/etc/rubrik-exporter/")
    v.AddConfigPath("$HOME/.rubrik-exporter")

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
    if err != nil {
        return nil, err
    }

    return &cfg, nil
}
