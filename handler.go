package main

import (
	"context"
	"net/http"
	"net/http/httputil"
	"sync/atomic"
	"time"
)

var serverPool *ServerPool

// One proxy for the app, that chooses backend dynamically for each request
var proxy *httputil.ReverseProxy

func initProxy() {
	proxy = &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			next, err := serverPool.GetNextValidPeer()
			if err != nil {
				// no backend alive. Mark in context that no backend is selected
				ctx := context.WithValue(req.Context(), "no-backend", false)
				req = req.WithContext(ctx)
				// request will fail and ErrorHandler will handle
				return
			}
			backend := serverPool.Backends[next]

			// increment connections
			atomic.AddInt64(&backend.CurrentConns, 1)

			// attach backend to request context for later decrement in ModifyResponse or ErrorHandler
			ctx := context.WithValue(req.Context(), "backend", backend)
			req = req.WithContext(ctx)

			req.URL.Scheme = backend.URL.Scheme
			req.URL.Host = backend.URL.Host
			req.Host = backend.URL.Host
		},

		ModifyResponse: func(resp *http.Response) error {
			// decrement connection count after response
			if b, ok := resp.Request.Context().Value("backend").(*Backend); ok {
				atomic.AddInt64(&b.CurrentConns, -1)
			}
			return nil
		},
		
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			// decrement connection count if request failed
			if b, ok := r.Context().Value("backend").(*Backend); ok {
				atomic.AddInt64(&b.CurrentConns, -1)
			}
			http.Error(w, "Backend unavailable", http.StatusServiceUnavailable)
		},
	}
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	// set a timeout context
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// clone the request with the timeout
	req := r.Clone(ctx)

	proxy.ServeHTTP(w, req)
}
