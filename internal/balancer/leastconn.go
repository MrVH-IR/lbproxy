package balancer

import "libproxy/internal/backend"

type LeastConn struct{ backends []*backend.Backend }

func NewLeastConn(backends []*backend.Backend) *LeastConn { return &LeastConn{backends: backends} }
func (lc *LeastConn) Next() *backend.Backend {
	best := lc.backends[0]
	for _, b := range lc.backends[1:] {
		if b.Active() < best.Active() {
			best = b
		}
	}
	return best
}
