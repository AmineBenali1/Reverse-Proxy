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
	if err := json.NewDecoder(r.Body).Decode(&newBackend); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	serverPool.Backends = append(serverPool.Backends, newBackend)
	w.Write([]byte("Backend added"))
}

func deleteBackend(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/backends" {
		http.Error(w, "Path not found", http.StatusNotFound)
		return
	}
	var backendToDelete Backend
	if err := json.NewDecoder(r.Body).Decode(&backendToDelete); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	for i, backend := range serverPool.Backends {
		if backend.URL == backendToDelete.URL {
			serverPool.Backends = append(serverPool.Backends[:i], serverPool.Backends[i+1:]...)
			w.Write([]byte("Backend Deleted Successfully"))
			return
		}
	}
	http.Error(w, "Backend not found", http.StatusNotFound)
}
