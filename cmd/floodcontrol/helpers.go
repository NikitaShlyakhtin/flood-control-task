package main

import (
	"flood-control/internal/floodcontrol"

	"github.com/redis/go-redis/v9"
)

func newLimiter(cfg config) floodcontrol.FloodControl {
	return floodcontrol.New(
		floodcontrol.Config{
			RPS:   float64(cfg.limiter.K / cfg.limiter.N),
			Burst: cfg.limiter.K,
		},
		openClient(cfg),
	)
}

func openClient(cfg config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.redis.Addr,
		Password: cfg.redis.Password,
		DB:       cfg.redis.DB,
	})
}
