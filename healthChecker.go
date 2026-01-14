package main

import (
	"fmt"
	"net/http"
	"time"
)

func PingBackend(b *Backend) {
	url := b.URL
	resp, err := http.Get(url.String())
	if err != nil {
		fmt.Println(url,"is DOWN:",err)
		b.mux.Lock()
		b.Alive = false
		b.CurrentConns = 0
		b.mux.Unlock()
		return
	}
	defer resp.Body.Close()

	b.mux.Lock()
	b.Alive = true
	b.mux.Unlock()
	fmt.Println(url,"is UP")
}

func startMonitoring(){
	ticker := time.NewTicker(proxyConfig.HealthCheckFreq)
	go func(){
		for range ticker.C{
			fmt.Println("Pinging backends...")
			for _ , b := range serverPool.Backends{
				go PingBackend(b)
			}
		}
	}()
}