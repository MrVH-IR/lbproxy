package backend

import (
	"fmt"
	"net/url"
	"sync/atomic"
)

type Backend struct {
	URL    *url.URL
	active atomic.Int64
}

func (b *Backend) Active() int64 { return b.active.Load() }
func (b *Backend) Start()        { b.active.Add(1) }
func (b *Backend) Done()         { b.active.Add(-1) }

func New(rawURL string) (*Backend, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("parse backend url %q: %w", rawURL, err)
	}
	if (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return nil, fmt.Errorf("backend url %q must look like http://host:port", rawURL)
	}
	return &Backend{URL: u}, nil
}
