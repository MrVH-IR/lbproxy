package balancer

import "libproxy/internal/backend"

type Balancer interface {
	Next() *backend.Backend
}
