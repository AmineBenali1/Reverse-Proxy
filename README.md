# Reverse-Proxy

I'll implement round-Robin first because i think it's easier, then I'll add least connections strategy and select between them dynamically

## 📌 TODO List

### 🔹 Core Architecture
- [x] Initialize Go module (`go mod init`)
- [x] Define core data structures (`Backend`, `ServerPool`, `ProxyConfig`)
- [x] Implement thread-safe server pool with mutexes / atomics
- [x] Handle case when no backend is available 

---

### 🔹 Load Balancing
- [x] Define `LoadBalancer` interface
- [x] Implement **Round-Robin** load balancing strategy
- [x] Implement **Least-Connections** load balancing strategy (at the end)
- [x] Ensure only healthy backends are selected
- [x] Allow dynamic strategy selection via configuration

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
- [x] Implement background health checker using goroutines
- [x] Use `time.Ticker` for periodic health checks
- [x] Ping backend servers to verify availability
- [x] Update backend alive status safely
- [x] Log backend state changes (UP / DOWN)

---

### 🔹 Admin API
- [x] Implement Admin API on a separate port
- [x] `GET /status` — show backend health and connection counts
- [x] `POST /backends` — add a new backend dynamically
- [x] `DELETE /backends` — remove an existing backend
- [x] Validate input and handle duplicate backend URLs
- [x] Return JSON responses

---

### 🔹 Configuration & Startup
- [x] Load proxy configuration from JSON file
- [x] Support command-line flag `--config`
- [x] Initialize server pool from config (optional empty start)
- [x] Start proxy server and admin API concurrently

---

### 🔹 Concurrency & Safety
- [x] Protect shared state using `sync.RWMutex`
- [x] Use `sync/atomic` for connection counters
- [x] Avoid race conditions under concurrent requests
- [ ] Verify code using Go race detector (`go run -race`)

---

### 🔹 Graceful Behavior
- [x] Handle client cancellation using `context.Context`
- [x] Cancel backend requests if client disconnects
- [x] Handle backend timeouts gracefully

---

### 🔹 Logging & Error Handling
- [ ] Log incoming requests and backend selection
- [x] Log health check failures and recoveries
- [x] Handle proxy errors and mark backends as unhealthy

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
