package main

import (
	"errors"
	"math"
	"strings"
	"sync"
	"sync/atomic"
	"net/url"
)

type ServerPool struct {
	Backends []*Backend `json:"backends"`
	Current  uint64     `json:"current"` // Used for Round-Robin
	mu sync.RWMutex
}

// Round-Robin

func (sp *ServerPool) GetNextValidPeer() (uint64,error) {
	for range sp.Backends{
		next := atomic.AddUint64(&sp.Current, 1) % uint64(len(sp.Backends))

		if sp.Backends[next].Alive {
			return next , nil 
		}
	}
	return math.MaxUint64 , errors.New("No available backend")
}

func (sp *ServerPool) AddBackend(backend *Backend) error{
	sp.mu.Lock()
	defer sp.mu.Unlock()

	for _, oldBackend := range sp.Backends{
		if strings.ToLower(oldBackend.URL.String()) == strings.ToLower(backend.URL.String()){
			return errors.New("Trying to add existing backend")
		}
	}
	
	sp.Backends = append(sp.Backends, backend)
	return nil
}

func (sp *ServerPool) SetBackendStatus(uri *url.URL, alive bool) error{
	sp.mu.Lock()
	defer sp.mu.Unlock()

	for _ , backend := range sp.Backends {
		if backend.URL == uri{
			backend.Alive = alive
			return nil 
		}
	}

	return errors.New("Backend not found")
}