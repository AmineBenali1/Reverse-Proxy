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
	if proxyConfig.Strategy != "round_robin" && proxyConfig.Strategy != "least_connections"{
		panic("Strategy not defined ! available strategies are round robin and least connections")
	}
	fmt.Println("proxy running on port", proxyConfig.Port)
	fmt.Println("Load balancing strategy:", proxyConfig.Strategy)
	fmt.Println("Health check frequency:", proxyConfig.HealthCheckFreq)
}
