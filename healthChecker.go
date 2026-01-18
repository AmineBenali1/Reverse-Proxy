package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func PingBackend(b *Backend) bool {
	url := b.URL
	resp, err := http.Get(url.String())
	if err != nil {
		fmt.Println(url, "is DOWN:", err)
		return false
	}
	defer resp.Body.Close()
	fmt.Println(url, "is UP")
	return true
}

func startMonitoring() {
	ticker := time.NewTicker(proxyConfig.HealthCheckFreq)
	go func() {
		for range ticker.C {
			fmt.Println("Pinging backends...")
			var wg sync.WaitGroup
			if serverPool == nil || len(serverPool.Backends) == 0 {
				fmt.Println("No backends, skipping health check.")
				continue
			}
			for _, b := range serverPool.Backends {
				wg.Add(1)
				go func(backend *Backend) {
					defer wg.Done()
					alive := PingBackend(backend)
					serverPool.SetBackendStatus(backend.URL, alive)
				}(b)
			}
			wg.Wait()
			fmt.Println("All backends pinged.")
		}
	}()
}
