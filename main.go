package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {

	proxyConfig = ProxyConfig{8080, "round-robin", 10 * time.Second}

	initProxy()
	startMonitoring()

	go StartAdminAPI()

	addr := fmt.Sprintf(":%d", proxyConfig.Port)

	http.HandleFunc("/", proxyHandler)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
