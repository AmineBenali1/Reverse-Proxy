package main

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"
)

var serverPool *ServerPool

// using NewSingleHostReverseProxy
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	// set context with timeout of 5s
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// get the next valid backend's id
	next, err := serverPool.GetNextValidPeer()
	if err != nil {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	backend := serverPool.Backends[next]
	// increment the current connections using atomic
	atomic.AddInt64(&backend.CurrentConns, 1)

	uri, err := url.Parse(backend.URL.String())
	if err != nil {
		http.Error(w, "Bad backend URL", http.StatusInternalServerError)
		return
	}

	// proxy per request (will change it to proxy per backend later)
	proxy := httputil.NewSingleHostReverseProxy(uri)

	// clone the context and Serve the request
	req := r.Clone(ctx)
	proxy.ServeHTTP(w, req)

	// decrement the current connections after finishing proccessing the request
	proxy.ModifyResponse = func(resp *http.Response) error {
		atomic.AddInt64(&backend.CurrentConns, -1)
		return nil
	}

	// decrement also if the request fails
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		atomic.AddInt64(&backend.CurrentConns, -1)
		http.Error(w, "Backend unavailable", http.StatusServiceUnavailable)
	}

}
