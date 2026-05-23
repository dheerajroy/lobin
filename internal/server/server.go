package server

import (
	"fmt"
	"net/http"

	"github.com/dheerajroy/lobin/internal/config"
	"github.com/dheerajroy/lobin/internal/logger"
	"github.com/dheerajroy/lobin/internal/strategy"
	"github.com/dheerajroy/lobin/internal/upstream"
)

type Server struct {
	upstreams []*upstream.Upstream
	strategy  strategy.Strategy
	config    *config.Config
}

func NewServer(upstreams []*upstream.Upstream, strategy strategy.Strategy, cfg *config.Config) *Server {
	return &Server{
		upstreams: upstreams,
		strategy:  strategy,
		config:    cfg,
	}
}

func (s *Server) getActiveUpstreams() []*upstream.Upstream {
	active := make([]*upstream.Upstream, 0)

	for _, u := range s.upstreams {
		if u.IsAlive() {
			active = append(active, u)
		}
	}
	return active
}

func (s *Server) loadBalancerHandler(w http.ResponseWriter, r *http.Request) {
	activeUpstreams := s.getActiveUpstreams()
	upstream, err := s.strategy.Select(activeUpstreams)
	if err != nil {
		logger.Error.Println("No active upstreams available")
		http.Error(
			w,
			"Service Unavailable",
			http.StatusServiceUnavailable,
		)
		return
	}

	upstream.IncrementConnections()
	defer upstream.DecrementConnections()

	upstream.Proxy.ServeHTTP(w, r)
}

func (s *Server) Serve() {

	http.HandleFunc("/", s.loadBalancerHandler)

	addr := fmt.Sprintf(
		":%d",
		s.config.Server.Port,
	)

	if s.config.TLS.Enabled {

		logger.Info.Printf("Load balancer listening on %s (TLS mode)\n", addr)

		err := http.ListenAndServeTLS(
			addr,
			s.config.TLS.CertFile,
			s.config.TLS.KeyFile,
			nil,
		)

		if err != nil {
			panic(err)
		}

		return
	}

	logger.Info.Printf("Load balancer listening on %s (HTTP mode)\n", addr)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
