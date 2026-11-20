package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Listen           string   `yaml:"listen"`
	Algorithm        string   `yaml:"algorithm"`
	Backends         []string `yaml:"backends"`
	LogDir           string   `yaml:"log_dir"`
	LogRetentionDays int      `yaml:"log_retention_days"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.Listen == "" {
		return errors.New("listen is required")
	}
	if len(c.Backends) == 0 {
		return errors.New("at least one backend is required")
	}
	if c.Algorithm == "" {
		c.Algorithm = "least_conn"
	}
	c.Algorithm = strings.ToLower(strings.TrimSpace(c.Algorithm))
	if c.Algorithm != "least_conn" && c.Algorithm != "round_robin" {
		return fmt.Errorf("unsupported algorithm %q (use least_conn or round_robin)", c.Algorithm)
	}
	if c.LogDir == "" {
		c.LogDir = "logs"
	}
	if c.LogRetentionDays == 0 {
		c.LogRetentionDays = 7
	}
	if c.LogRetentionDays < 1 {
		return errors.New("log_retention_days must be at least 1")
	}
	return nil
}
