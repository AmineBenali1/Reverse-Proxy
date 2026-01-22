package main

import (
	"encoding/json"
	"errors"
	"math"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

type ServerPool struct {
	Backends []*Backend `json:"backends"`
	Current  uint64     `json:"current"` // Used for Round-Robin
	mu       sync.RWMutex
}

// Round-Robin

func (sp *ServerPool) GetNextValidPeer() (uint64, error) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	for range sp.Backends {
		next := atomic.AddUint64(&sp.Current, 1) % uint64(len(sp.Backends))

		if sp.Backends[next].Alive {
			return next, nil
		}
	}
	return math.MaxUint64, errors.New("No available backend")
}

func (sp *ServerPool) AddBackend(backend *Backend) error {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	for _, oldBackend := range sp.Backends {
		if strings.EqualFold(oldBackend.URL.String(), backend.URL.String()) {
			return errors.New("Trying to add existing backend")
		}
	}

	sp.Backends = append(sp.Backends, backend)
	return nil
}

func (sp *ServerPool) SetBackendStatus(uri *url.URL, alive bool) error {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	for _, backend := range sp.Backends {
		if backend.URL == uri {
			backend.Alive = alive
			if !alive {
				backend.CurrentConns = 0
			}
			return nil
		}
	}

	return errors.New("Backend not found")
}

func LoadDefinedBackends(backendConfigPath string) {
	file, err := os.Open(backendConfigPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Temporary structs with URL as string
	type tmpBackend struct {
		URL          string `json:"url"`
		Alive        bool   `json:"alive"`
		CurrentConns int64  `json:"current_connections"`
	}

	type tmpServerPool struct {
		Backends []tmpBackend `json:"backends"`
		Current  uint64       `json:"current"`
	}

	var tmpPool tmpServerPool
	if err := json.NewDecoder(file).Decode(&tmpPool); err != nil {
		panic(err)
	}

	// convert to actual server pool
	serverPool = &ServerPool{Current: tmpPool.Current}
	for _, b := range tmpPool.Backends {
		ParsedUrl, err := url.Parse(b.URL)
		if err != nil {
			panic(err)
		}
		serverPool.Backends = append(serverPool.Backends, &Backend{
			URL:          ParsedUrl,
			Alive:        b.Alive,
			CurrentConns: b.CurrentConns,
		})
	}
}
