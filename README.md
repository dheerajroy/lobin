# Lobin - Lightweight Load Balancer

<div align="center">
  <img src="./assets/lobin.png" alt="Lobin Logo" width="300" />
</div>

Lobin is a fast, lightweight HTTP load balancer written in Go. It distributes incoming requests across multiple backend servers with built-in health checking, multiple load balancing strategies, and optional TLS termination.

**⚠️ PRODUCTION WARNING** - See [Production Considerations](#production-considerations) below.

## Features

- **Multiple Load Balancing Strategies**
  - Round Robin - Distribute requests evenly across backends
  - Weighted Round Robin - Distribute based on backend capacity/priority
  - Least Connections - Route to backend with fewest active connections

- **Automatic Health Checking**
  - Periodic health checks on all upstreams
  - Automatic recovery when backends become healthy again
  - Configurable check intervals and timeouts

- **TLS/HTTPS Support**
  - Optional TLS termination for client connections
  - Configurable certificate and key paths
  - Automatic HTTP fallback for non-TLS deployments

- **Structured Logging**
  - Multiple log levels (INFO, DEBUG, ERROR)
  - Clear visibility into routing decisions and upstream health
  - Configurable logging level via config file

- **Thread-Safe Operation**
  - Concurrent request handling with proper synchronization
  - Safe state management for upstreams and connections

## Installation

### Build from Source

```bash
# Clone the repository
git clone https://github.com/dheerajroy/lobin.git
cd lobin

# Build the binary
go build -o lobin ./cmd/lobin

# Run with config file
./lobin -config ./config/config.yaml
```

### Requirements
- Go 1.16 or higher
- Standard Go libraries only (no external dependencies)

## Configuration

Lobin uses YAML for configuration. A sample config is provided at `config/config.yaml`.

### Basic Configuration Structure

```yaml
server:
  port: 8000                    # Listen port

strategy: round_robin           # Load balancing strategy

health_check:
  interval: 10s                # Health check frequency
  timeout: 5s                  # Health check timeout

upstreams:                      # Backend servers
  - url: http://localhost:8001
    weight: 1
  - url: http://localhost:8002
    weight: 1

tls:
  enabled: false               # Enable HTTPS
  # cert_file: ./server.crt
  # key_file: ./server.key

logging:
  level: info                  # Log level: info, debug, error
```

### Configuration Options

#### Server

```yaml
server:
  port: 8000                   # Port to listen on (default: 8000)
```

#### Strategy

Choose one of:
- `round_robin` - Even distribution
- `weighted_round_robin` - Distribution based on `weight` values
- `least_connections` - Route to upstream with fewest active connections

```yaml
strategy: round_robin
```

#### Health Check

```yaml
health_check:
  interval: 10s                # How often to check backends
  timeout: 5s                  # Timeout for each health check
```

#### Upstreams

Define your backend servers:

```yaml
upstreams:
  - url: http://localhost:8001
    health_check_path: /health  # Endpoint to check
    weight: 1                   # Weight (optional, default: 1)
  - url: http://localhost:8002
    health_check_path: /health
    weight: 2                   # This backend gets 2x more traffic
```

**Weight Configuration:**
- `weight` is optional - defaults to 1 if not specified
- Only used by `weighted_round_robin` strategy
- Ignored by `round_robin` and `least_connections` strategies
- Higher weight = more traffic (proportional distribution)

#### TLS Configuration

**HTTP Mode** (default):
```yaml
tls:
  enabled: false
```

**HTTPS Mode**:
```yaml
tls:
  enabled: true
  cert_file: ./certs/server.crt
  key_file: ./certs/server.key
```

Generate self-signed certificates for testing:
```bash
openssl req -x509 -newkey rsa:4096 \
  -keyout server.key \
  -out server.crt \
  -days 365 \
  -nodes
```

#### Logging

```yaml
logging:
  level: info                  # Options: info, debug, error
```

## Usage

### Quick Start

1. Configure your upstreams in `config/config.yaml`
2. Run the load balancer:
   ```bash
   ./lobin -config ./config/config.yaml
   ```
3. Send requests to the load balancer:
   ```bash
   curl http://localhost:8000/
   ```

### Example: Round Robin with 2 Backends

**config.yaml:**
```yaml
server:
  port: 8000

strategy: round_robin

health_check:
  interval: 10s
  timeout: 5s

upstreams:
  - url: http://localhost:3001
    health_check_path: /health
    weight: 1
  - url: http://localhost:3002
    health_check_path: /health
    weight: 1

logging:
  level: info
```

**Run:**
```bash
# Terminal 1: Start backend 1
python -m http.server 3001

# Terminal 2: Start backend 2
python -m http.server 3002

# Terminal 3: Start Lobin
./lobin -config ./config/config.yaml

# Terminal 4: Send requests
curl http://localhost:8000/
```

### Example: Weighted Load Balancing

Route 2x traffic to a more powerful server:

```yaml
upstreams:
  - url: http://powerful-server:8000
    health_check_path: /health
    weight: 2              # Gets 2x more requests
  - url: http://standard-server:8000
    health_check_path: /health
    weight: 1              # Gets standard amount
```

With 3 requests: Powerful gets 2, Standard gets 1.

## Architecture

```
┌─────────────┐
│   Clients   │
└──────┬──────┘
       │
       ▼
┌──────────────────────────────┐
│   Lobin Load Balancer        │
├──────────────────────────────┤
│ • Request Handler            │
│ • Strategy Selector          │
│ • Connection Tracker         │
└──────┬───────────────────────┘
       │
   ┌───┴────┬──────────┬────────┐
   ▼        ▼          ▼        ▼
┌────────┬────────┬────────┬────────┐
│Backend │Backend │Backend │Backend │
│   1    │   2    │   3    │   4    │
└────────┴────────┴────────┴────────┘
   │        │          │        │
   └────┬───┴──────┬───┴────┬───┘
        ▼          ▼        ▼
   ┌─────────────────────────────┐
   │  Health Check Goroutines    │
   │  (Periodic verification)    │
   └─────────────────────────────┘
```

## How It Works

1. **Request Reception** - Client request arrives at Lobin
2. **Strategy Selection** - Based on configured strategy, an upstream is selected
3. **Health Check** - Only healthy upstreams are considered
4. **Connection Tracking** - Connection count is incremented for chosen upstream
5. **Proxying** - Request is forwarded to selected upstream
6. **Response** - Response is sent back to client
7. **Cleanup** - Connection count is decremented

## Logging Output

When running with `level: debug`, you'll see:

```
[INFO] Logger initialized with level: debug
[INFO] Loaded 2 upstream(s) from config
[DEBUG] Starting health check for: http://localhost:8001
[DEBUG] Starting health check for: http://localhost:8002
[INFO] Load balancer listening on :8000 (HTTP mode)
```

Errors appear as:
```
[ERROR] Health check failed for http://localhost:8001: connection refused
[ERROR] Upstream marked as inactive: http://localhost:8001
[INFO] Upstream recovered and marked as active: http://localhost:8001
[ERROR] No active upstreams available
```

## Production Considerations

⚠️ **WARNING: This is a prototype/learning project** and should not be used in production environments without significant hardening and testing.

This load balancer is suitable for:
- Learning and understanding load balancing concepts
- Development and testing environments
- Small-scale internal use cases


## License

See LICENSE file for details.
