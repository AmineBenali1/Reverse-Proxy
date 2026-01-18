package main

import (
	"encoding/json"
	"net/http"
)

func getStatus(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/status" {
		http.Error(w, "Path not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var activeBackends int
	backends := make([]map[string]interface{}, 0, len(serverPool.Backends))

	for _, backend := range serverPool.Backends {
		details := map[string]interface{}{
			"url":                 backend.URL,
			"alive":               backend.Alive,
			"current_connections": backend.CurrentConns,
		}
		backends = append(backends, details)
		if backend.Alive {
			activeBackends++
		}
	}

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
	if r.URL.Path != "/backends" {
		http.Error(w, "Path not found", http.StatusNotFound)
		return
	}
	var newBackend *Backend
	json.NewDecoder(r.Body).Decode(&newBackend)
	serverPool.Backends = append(serverPool.Backends, newBackend)
	w.Write([]byte("Backend added"))
}


