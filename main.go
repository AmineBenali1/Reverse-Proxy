package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

func main() {

	proxyConfig = ProxyConfig{8080, "round-robin", 10 * time.Second}
	url1, _ := url.Parse("http://localhost:8081")
	url2, _ := url.Parse("http://localhost:8082")
	url3, _ := url.Parse("http://localhost:8083")
	backend1 := &Backend{URL: url1, Alive: true, CurrentConns: 0}
	backend2 := &Backend{URL: url2, Alive: true, CurrentConns: 0}
	backend3 := &Backend{URL: url3, Alive: true, CurrentConns: 0}
	backends := []*Backend{}
	backends = append(backends, backend1, backend2, backend3)
	serverPool = &ServerPool{Backends: backends}

	initProxy()
	startMonitoring()

	addr := fmt.Sprintf(":%d", proxyConfig.Port)

	http.HandleFunc("/", proxyHandler)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
