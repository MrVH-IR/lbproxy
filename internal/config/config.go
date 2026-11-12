package config

import (
	"errors"
	"fmt"
	"os"
)

type Config struct {
	Listen    string   `yaml:"listen"`
	Algorithm string   `yaml:"algorithm"`
	Backends  []string `yaml:"backends"`
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
	// TODO (تمرین ۱): دو چک دیگه
	// - اگه len(c.Backends) == 0 خطا بده
	// - اگه c.Algorithm مقدارش "round_robin" نبود خطا بده
	//   (با fmt.Errorf و %q اسم مقدار نامعتبر رو توی پیام بذار)
	return nil
}
