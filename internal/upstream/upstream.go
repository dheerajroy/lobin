package upstream

import (
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/dheerajroy/lobin/internal/logger"
)

type Upstream struct {
	URL                 *url.URL
	HealthCheckPath     string
	Alive               bool
	Proxy               *httputil.ReverseProxy
	Weight              int
	currentWeight       int
	mu                  sync.RWMutex
	connections         int
	HealthCheckInterval time.Duration
	HealthCheckTimeout  time.Duration
}

func NewUpstream(rawURL string, healthCheckPath string, weight int, healthCheckInterval, healthCheckTimeout time.Duration) *Upstream {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}

	u := &Upstream{
		URL:                 parsedURL,
		HealthCheckPath:     healthCheckPath,
		Alive:               true,
		Proxy:               httputil.NewSingleHostReverseProxy(parsedURL),
		Weight:              weight,
		currentWeight:       0,
		connections:         0,
		HealthCheckInterval: healthCheckInterval,
		HealthCheckTimeout:  healthCheckTimeout,
	}
	logger.Debug.Printf("Created upstream: %s (weight: %d)\n", parsedURL.String(), weight)
	return u
}

func (u *Upstream) IsAlive() bool {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.Alive
}

func (u *Upstream) GetConnections() int {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.connections
}

func (u *Upstream) GetCurrentWeight() int {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.currentWeight
}

func (u *Upstream) SetCurrentWeight(weight int) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.currentWeight = weight
}

func (u *Upstream) IncrementConnections() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.connections++
}

func (u *Upstream) DecrementConnections() {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.connections > 0 {
		u.connections--
	}
}

func (u *Upstream) SetInactive() {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.Alive {
		logger.Error.Printf("Upstream marked as inactive: %s\n", u.URL.String())
		u.Alive = false
	}
}

func (u *Upstream) SetActive() {
	u.mu.Lock()
	defer u.mu.Unlock()
	if !u.Alive {
		logger.Info.Printf("Upstream recovered and marked as active: %s\n", u.URL.String())
		u.Alive = true
	}
}
