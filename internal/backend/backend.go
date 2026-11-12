package backend

import (
	"fmt"
	"net/url"
)

type Backend struct {
	URL *url.URL
}

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
