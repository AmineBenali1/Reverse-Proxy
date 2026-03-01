
# Reverse Proxy

A high-performance, concurrent HTTP reverse proxy and load balancer written in Go. Supports dynamic backend management, health checks, and multiple load balancing strategies.

## Features

- **Load Balancing:** Round-Robin and Least-Connections strategies (configurable)
- **Dynamic Backend Management:** Add or remove backends at runtime via Admin API
- **Health Monitoring:** Periodic health checks to ensure only healthy backends receive traffic
- **Graceful Error Handling:** Handles backend failures and client cancellations
- **Concurrency:** Thread-safe backend pool and connection counters
- **Configuration:** JSON-based config for easy setup

## Architecture Overview

The proxy consists of several core components:

- **Proxy Server:** Receives client requests and forwards them to healthy backends.
- **Server Pool:** Maintains a list of backend servers, their health, and connection counts.
- **Load Balancer:** Selects the next backend using the configured strategy.
- **Health Checker:** Periodically pings backends to update their status.
- **Admin API:** Allows runtime management and monitoring of backends.


## Load Balancing Strategies

- **Round-Robin:**
	- Each incoming request is forwarded to the next available backend in a circular order.
	- Simple and effective for evenly distributed workloads.

- **Least-Connections:**
	- The backend with the fewest active connections is selected for each new request.
	- Useful for workloads with varying request durations.

You can switch strategies by editing the `strategy` field in `config.json`.

## Health Check Mechanism

- The proxy periodically sends HTTP GET requests to each backend.
- If a backend responds successfully, it is marked as healthy (`alive: true`).
- Unhealthy backends are automatically excluded from load balancing until they recover.
- Health check frequency is set in nanoseconds in `config.json` (e.g., `10000000000` for 10 seconds).

## Project Structure

- `main.go` — Entry point, starts proxy and admin API
- `handler.go` — Reverse proxy logic
- `serverPool.go` — Backend pool and load balancing
- `adminAPI.go` — Admin API for backend management and status
- `healthChecker.go` — Health check routines
- `proxyConfig.go` — Loads and validates configuration
- `backends/` — Example backend servers

## Getting Started

### Prerequisites

- Go 1.25+ installed

### Installation & Running

1. **Clone the repository:**
	 ```bash
	 git clone https://github.com/AmineBenali1/Reverse-Proxy.git
	 cd Reverse-Proxy
	 ```

2. **Build and run the proxy:**
	 ```bash
	 go run main.go
	 ```
	 - Use `--config` to specify a custom config file.
	 - Use `--backend-config` to load initial backends from a file.

3. **Start example backends:**
	 ```bash
	 go run backends/backend1.go
	 go run backends/backend2.go
	 go run backends/backend3.go
	 ```

## Configuration

Edit `config.json` to set the proxy port, load balancing strategy, and health check frequency:
```json
{
	"port": 8080,
	"strategy": "round_robin", // or "least_connections"
	"health_check_frequency": 10000000000
}
```

#### Example `backendConfig.json`
```json
{
	"backends": [
		{
			"url": "http://localhost:8082",
			"alive": false,
			"current_connections": 0
		}
	],
	"current": 0
}
```

## Admin API

Runs on port `8081` by default.

### Endpoints

- `GET /status` — Returns status of all backends.
- `POST /backends` — Add a backend. Body: `{ "url": "http://localhost:8084" }`
- `DELETE /backends` — Remove a backend. Body: `{ "url": "http://localhost:8084" }`

### Example JSON Response

**GET /status**
```json
{
	"total_backends": 3,
	"active_backends": 2,
	"backends": [
		{
			"url": "http://localhost:8082",
			"alive": true,
			"current_connections": 0
		},
		{
			"url": "http://localhost:8083",
			"alive": true,
			"current_connections": 1
		},
		{
			"url": "http://localhost:8084",
			"alive": false,
			"current_connections": 0
		}
	]
}
```

### Example Commands

#### Add a Backend

- **cURL (Linux/macOS/Windows):**
	```bash
	curl -X POST http://localhost:8081/backends -H "Content-Type: application/json" -d '{"url": "http://localhost:8084"}'
	```

- **Windows PowerShell:**
	```powershell
	Invoke-RestMethod -Uri "http://localhost:8081/backends" -Method Post -ContentType "application/json" -Body '{"url": "http://localhost:8084"}'
	```

#### Delete a Backend

- **cURL (Linux/macOS/Windows):**
	```bash
	curl -X DELETE http://localhost:8081/backends -H "Content-Type: application/json" -d '{"url": "http://localhost:8084"}'
	```

- **Windows PowerShell:**
	```powershell
	Invoke-RestMethod -Uri "http://localhost:8081/backends" -Method Delete -ContentType "application/json" -Body '{"url": "http://localhost:8084"}'
	```

#### Get Backend Status

- **cURL (Linux/macOS/Windows):**
	```bash
	curl -X GET http://localhost:8081/status
	```

- **Windows PowerShell:**
	```powershell
	Invoke-RestMethod -Uri "http://localhost:8081/status" -Method Get
	```

## Example Backend

Start a backend server (e.g., on port 8084):
```bash
go run backends/backend1.go
```

## Troubleshooting & FAQ

**Q: The proxy does not forward requests to a backend.**
- Check if the backend is running and healthy (see `/status`).
- Ensure the backend URL is correct and accessible.

**Q: How do I change the load balancing strategy?**
- Edit the `strategy` field in `config.json` to `round_robin` or `least_connections` and restart the proxy.

**Q: How do I add or remove backends at runtime?**
- Use the Admin API endpoints as shown above.

**Q: How do I run the proxy on a different port?**
- Change the `port` value in `config.json`.

**Q: How do I increase health check frequency?**
- Lower the `health_check_frequency` value (in nanoseconds) in `config.json`.

