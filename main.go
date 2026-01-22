package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// flags
	configPath := flag.String("config", "config.json", "")
	backendConfigPath := flag.String("backend-config", "", "")
	flag.Parse()

	LoadProxyConfig(*configPath)
	if *backendConfigPath != "" {
		LoadDefinedBackends(*backendConfigPath)
	}

	initProxy()
	startMonitoring()

	go StartAdminAPI()

	addr := fmt.Sprintf(":%d", proxyConfig.Port)

	http.HandleFunc("/", proxyHandler)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
