package strategy

import (
	"github.com/dheerajroy/lobin/internal/upstream"
	"sync"
	"sync/atomic"
)

type RoundRobin struct {
	index uint64
	mu    sync.Mutex
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{}
}

func (rr *RoundRobin) Select(activeUpstreams []*upstream.Upstream) (*upstream.Upstream, error) {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	n := len(activeUpstreams)
	if n == 0 {
		return nil, ErrNoUpstreams
	}

	i := atomic.AddUint64(&rr.index, 1) - 1
	idx := i % uint64(n)
	return activeUpstreams[idx], nil
}
