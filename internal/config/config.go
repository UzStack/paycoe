package config

import (
	"os"
	"strconv"
)

type Config struct {
	Workers int
}

func NewConfig() (*Config, error) {
	workers_str := os.Getenv("WORKERS")
	if workers_str == "" {
		workers_str = "10"
	}
	workers, err := strconv.Atoi(workers_str)
	if err != nil {
		return nil, err
	}
	return &Config{
		Workers: workers,
	}, nil
}
