package health

import (
	"context"
	"net/http"
	"time"

	"github.com/dheerajroy/lobin/internal/logger"
	"github.com/dheerajroy/lobin/internal/upstream"
)

func Start(ctx context.Context, u *upstream.Upstream) {
	ticker := time.NewTicker(u.HealthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			check(u)

		case <-ctx.Done():
			return
		}
	}
}

func check(u *upstream.Upstream) {
	client := &http.Client{
		Timeout: u.HealthCheckTimeout,
	}

	resp, err := client.Get(u.URL.String() + u.HealthCheckPath)
	if err != nil {
		logger.Error.Printf("Health check failed for %s: %v\n", u.URL.String(), err)
		u.SetInactive()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		u.SetActive()
	} else {
		logger.Error.Printf("Health check failed for %s: status code %d\n", u.URL.String(), resp.StatusCode)
		u.SetInactive()
	}
}
