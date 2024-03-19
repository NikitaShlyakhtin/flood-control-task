package main

import (
	"context"
	"flag"
	"flood-control/internal/floodcontrol"
	"fmt"
	"time"
)

type config struct {
	limiter struct {
		N int // Rate limit window in seconds
		K int // Maximum number of requests per N seconds
	}
	redis struct {
		Addr     string
		Password string
		DB       int
	}
}

type application struct {
	config config
	fc     floodcontrol.FloodControl
}

func main() {
	var cfg config

	flag.IntVar(&cfg.limiter.N, "limiter.n", 2, "Rate limit window in seconds")
	flag.IntVar(&cfg.limiter.K, "limiter.k", 4, "Maximum number of requests per N seconds")

	flag.StringVar(&cfg.redis.Addr, "redis.addr", "localhost:6379", "Redis address")
	flag.StringVar(&cfg.redis.Password, "redis.password", "", "Redis password")
	flag.IntVar(&cfg.redis.DB, "redis.db", 0, "Redis database")

	flag.Parse()

	app := &application{
		config: cfg,
		fc:     newLimiter(cfg),
	}

	// Simulate multiple requests
	for i := 0; i < 20; i++ {
		// Perform rate limiting check for user ID 1
		result, err := app.fc.Check(context.Background(), 1)
		if err != nil {
			fmt.Printf("Error performing rate limiting check: %v\n", err)
			continue
		}

		if result {
			fmt.Println("Request allowed")
		} else {
			fmt.Println("Request denied")
		}

		// Sleep for a short duration to simulate requests being made over time
		time.Sleep(100 * time.Millisecond)
	}
}
