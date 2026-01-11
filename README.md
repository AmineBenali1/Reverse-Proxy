# Reverse-Proxy

I'll implement round-Robin first because i think it's easier, then I'll add least connections strategy and select between them dynamically

## 📌 TODO List

### 🔹 Core Architecture
- [x] Initialize Go module (`go mod init`)
- [x] Define core data structures (`Backend`, `ServerPool`, `ProxyConfig`)
- [x] Implement thread-safe server pool with mutexes / atomics
- [ ] Handle case when no backend is available (return `503 Service Unavailable`)

---

### 🔹 Load Balancing
- [x] Define `LoadBalancer` interface
- [x] Implement **Round-Robin** load balancing strategy
- [ ] Implement **Least-Connections** load balancing strategy (at the end)
- [x] Ensure only healthy backends are selected
- [ ] Allow dynamic strategy selection via configuration

---

### 🔹 Reverse Proxy
- [x] Implement HTTP proxy handler
- [x] Integrate `httputil.NewSingleHostReverseProxy`
- [x] Forward client requests to selected backend
- [x] Propagate request context to backend
- [x] Increment and decrement backend connection counters
- [x] Implement custom error handler to detect backend failures

---

### 🔹 Health Monitoring
- [ ] Implement background health checker using goroutines
- [ ] Use `time.Ticker` for periodic health checks
- [ ] Ping backend servers to verify availability
- [ ] Update backend alive status safely
- [ ] Log backend state changes (UP / DOWN)

---

### 🔹 Admin API
- [ ] Implement Admin API on a separate port
- [ ] `GET /status` — show backend health and connection counts
- [ ] `POST /backends` — add a new backend dynamically
- [ ] `DELETE /backends` — remove an existing backend
- [ ] Validate input and handle duplicate backend URLs
- [ ] Return JSON responses

---

### 🔹 Configuration & Startup
- [ ] Load proxy configuration from JSON file
- [ ] Support command-line flag `--config`
- [ ] Initialize server pool from config (optional empty start)
- [ ] Start proxy server and admin API concurrently

---

### 🔹 Concurrency & Safety
- [ ] Protect shared state using `sync.RWMutex`
- [ ] Use `sync/atomic` for connection counters
- [ ] Avoid race conditions under concurrent requests
- [ ] Verify code using Go race detector (`go run -race`)

---

### 🔹 Graceful Behavior
- [ ] Handle client cancellation using `context.Context`
- [ ] Cancel backend requests if client disconnects
- [ ] Handle backend timeouts gracefully

---

### 🔹 Logging & Error Handling
- [ ] Log incoming requests and backend selection
- [ ] Log health check failures and recoveries
- [ ] Handle proxy errors and mark backends as unhealthy

---

### 🔹 Documentation
- [ ] Write project overview in README
- [ ] Document load-balancing strategies
- [ ] Document Admin API endpoints with examples
- [ ] Provide instructions to run the project
- [ ] Add example `config.json`

---

### 🔹 Optional Enhancements
- [ ] Sticky sessions (IP / Cookie based)
- [ ] Weighted load balancing
- [ ] HTTPS / TLS support
- [ ] Backend persistence to disk
- [ ] Metrics (request count, latency)
