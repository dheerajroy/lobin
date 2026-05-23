package strategy

import (
	"github.com/dheerajroy/lobin/internal/upstream"
)

type LeastConnections struct {}

func NewLeastConnections() *LeastConnections {
	return &LeastConnections{}
}

func (lc *LeastConnections) Select(activeUpstreams []*upstream.Upstream) (*upstream.Upstream, error) {
	n := len(activeUpstreams)
	if n == 0 {
		return nil, ErrNoUpstreams
	}

	var best *upstream.Upstream
	for _, up := range activeUpstreams {
		if best == nil || up.GetConnections() < best.GetConnections() {
			best = up
		}
	}
	return best, nil
}
