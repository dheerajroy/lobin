package strategy

import (
	"github.com/dheerajroy/lobin/internal/upstream"
	"sync"
)

type WeightedRoundRobin struct {
	mu sync.Mutex
}

func NewWeightedRoundRobin() *WeightedRoundRobin {
	return &WeightedRoundRobin{}
}

func (wrr *WeightedRoundRobin) Select(activeUpstreams []*upstream.Upstream) (*upstream.Upstream, error) {
	wrr.mu.Lock()
	defer wrr.mu.Unlock()

	n := len(activeUpstreams)
	if n == 0 {
		return nil, ErrNoUpstreams
	}

	var best *upstream.Upstream
	totalWeight := 0

	for _, up := range activeUpstreams {
		weight := up.Weight
		currentWeight := up.GetCurrentWeight() + weight

		up.SetCurrentWeight(currentWeight)

		totalWeight += weight

		if best == nil || currentWeight > best.GetCurrentWeight() {
			best = up
		}
	}

	best.SetCurrentWeight(
		best.GetCurrentWeight() - totalWeight,
	)
	return best, nil
}
