package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

var proxyConfig ProxyConfig

type ProxyConfig struct {
	Port            int           `json:"port"`
	Strategy        string        `json:"strategy"` // e.g., "round-robin" or "least-conn"
	HealthCheckFreq time.Duration `json:"health_check_frequency"`
}

func LoadProxyConfig(configPath string) {
	file, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&proxyConfig); err != nil {
		panic(err)
	}

	fmt.Println("proxy running on port", proxyConfig.Port)
	fmt.Println("Load balancing strategy:", proxyConfig.Strategy)
	fmt.Println("Health check frequency:", proxyConfig.HealthCheckFreq)
}
