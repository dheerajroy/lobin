package strategy

import (
	"errors"
)

var ErrNoUpstreams = errors.New("no active upstreams available")
