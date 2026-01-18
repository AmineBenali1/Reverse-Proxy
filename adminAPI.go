package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

func getStatus(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if serverPool == nil {
		http.Error(w, "Server pool not initialized", http.StatusServiceUnavailable)
		return
	}
	var activeBackends int
	backends := make([]map[string]interface{}, 0, len(serverPool.Backends))
	serverPool.mu.RLock()
	for _, backend := range serverPool.Backends {
		details := map[string]interface{}{
			"url":                 backend.URL.String(),
			"alive":               backend.Alive,
			"current_connections": backend.CurrentConns,
		}
		backends = append(backends, details)
		if backend.Alive {
			activeBackends++
		}
	}
	serverPool.mu.RUnlock()
	response := map[string]interface{}{
		"total_backends":  len(serverPool.Backends),
		"active_backends": activeBackends,
		"backends":        backends,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error during operation", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}

func addBackend(w http.ResponseWriter, r *http.Request) {
	var newBackend *Backend
	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	parsedURL, err := url.Parse(req.URL)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	newBackend = &Backend{
		URL:          parsedURL,
		Alive:        false,
		CurrentConns: 0,
	}
	if serverPool == nil {
		serverPool = &ServerPool{}
	}
	if err := serverPool.AddBackend(newBackend); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.Write([]byte("Backend added"))
	}
}

func deleteBackend(w http.ResponseWriter, r *http.Request) {
	if serverPool == nil {
		http.Error(w, "Server pool not initialized", http.StatusServiceUnavailable)
		return
	}
	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	parsedURL, err := url.Parse(req.URL)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	for i, backend := range serverPool.Backends {
		if backend.URL.String() == parsedURL.String() {
			serverPool.mu.Lock()
			serverPool.Backends = append(serverPool.Backends[:i], serverPool.Backends[i+1:]...)
			serverPool.mu.Unlock()
			w.Write([]byte("Backend Deleted Successfully"))
			return
		}
	}
	http.Error(w, "Backend not found", http.StatusNotFound)
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/backends":
		switch r.Method {
		case http.MethodPost:
			addBackend(w, r)
		case http.MethodDelete:
			deleteBackend(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	case "/status":
		switch r.Method {
		case http.MethodGet:
			getStatus(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	case "/":
		w.Write([]byte("Debugging..."))
	default:
		http.Error(w, "Path not found", http.StatusNotFound)
	}
}

func StartAdminAPI() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Fatal(err)
	}
}
