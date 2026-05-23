package upstream

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"github.com/dheerajroy/lobin/internal/config"
)

func BuildFromConfig(
	cfg *config.Config,
) []*Upstream {

	var upstreams []*Upstream

	for _, up := range cfg.Upstreams {

		scheme := "http"

		if cfg.BackendTLS.Enabled {
			scheme = "https"
		}

		rawURL := fmt.Sprintf(
			"%s://%s",
			scheme,
			up.URL,
		)

		parsedURL, err := url.Parse(rawURL)
		if err != nil {
			panic(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(parsedURL)

		if cfg.BackendTLS.Enabled {

			transport := &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: cfg.BackendTLS.InsecureSkipVerify,
				},
			}

			proxy.Transport = transport
		}

		u := NewUpstream(
			rawURL,
			up.HealthCheckPath,
			up.Weight,
			cfg.HealthCheck.Interval,
			cfg.HealthCheck.Timeout,
		)

		u.Proxy = proxy

		upstreams = append(upstreams, u)
	}

	return upstreams
}