package hasher

import "fmt"

type Cost int

const (
	MinCost     Cost = 4
	MaxCost     Cost = 31
	DefaultCost Cost = 10
)

type Config struct {
	Cost Cost `env:"HASHER_COST" default:"10"`
}

func (cfg *Config) validate() error {
	if cfg.Cost < MinCost || cfg.Cost > MaxCost {
		return fmt.Errorf("cost must be between %d and %d", MinCost, MaxCost)
	}
	return nil
}
