package strategy

import (
	"github.com/dheerajroy/lobin/internal/upstream"
)

type Strategy interface {
	Select(activeUpstreams []*upstream.Upstream) (*upstream.Upstream, error)
}
