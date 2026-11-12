package balancer

import (
	"sync/atomic"

	"libproxy/internal/backend"
)

type RoundRobin struct {
	backends []*backend.Backend
	counter  atomic.Uint64
}

func NewRoundRobin(backends []*backend.Backend) *RoundRobin {
	return &RoundRobin{backends: backends}
}

func (rr *RoundRobin) Next() *backend.Backend {
	n := rr.counter.Add(1) - 1
	return rr.backends[n%uint64(len(rr.backends))]
}
