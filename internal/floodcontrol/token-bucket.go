package floodcontrol

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenBucket struct {
	config Config
	client *redis.Client
}

func New(config Config, client *redis.Client) *TokenBucket {
	return &TokenBucket{
		config: config,
		client: client,
	}
}

type client struct {
	Tokens     float64   `json:"tokens"`
	LastRefill time.Time `json:"last_refill"`
}

func (tb *TokenBucket) Check(ctx context.Context, userID int64) (bool, error) {
	key := fmt.Sprint(userID)

	// Start Redis transaction
	tx := tb.client.TxPipeline()

	// Get and update token bucket data within the transaction
	cl, err := tb.getClientData(ctx, tx, key)
	if err != nil {
		return false, err
	}

	result := tb.updateTokenBucket(ctx, tx, key, cl)

	// Execute the transaction
	_, err = tx.Exec(ctx)
	if err != nil {
		return false, err
	}

	return result, nil
}

func (tb *TokenBucket) getClientData(ctx context.Context, tx redis.Pipeliner, key string) (*client, error) {
	// Get token bucket data from Redis
	getCmd := tx.Get(ctx, key)
	_, err := tx.Exec(ctx)

	if err == redis.Nil {
		// If the key doesn't exist, create a new token bucket
		return &client{
			Tokens:     float64(tb.config.Burst),
			LastRefill: time.Now(),
		}, nil
	}

	if err != nil {
		return nil, err
	}

	// Unmarshal token bucket data
	var cl client

	data, err := getCmd.Bytes()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &cl); err != nil {
		return nil, err
	}

	return &cl, nil
}

func (tb *TokenBucket) updateTokenBucket(ctx context.Context, tx redis.Pipeliner, key string, cl *client) bool {
	tb.refillTokens(cl)
	result := tb.consumeToken(cl)

	// Store the updated client data back in Redis within the transaction
	data, err := json.Marshal(cl)
	if err != nil {
		return false
	}

	tx.Set(ctx, key, data, 0)

	return result
}

func (tb *TokenBucket) refillTokens(cl *client) {
	// Calculate the number of tokens that should have been refilled since the last refill
	now := time.Now()
	elapsed := now.Sub(cl.LastRefill)
	refillTokens := float64(elapsed.Seconds() * tb.config.RPS)

	// Add the refilled tokens to the client's token bucket, up to the maximum burst size
	cl.Tokens += refillTokens
	if cl.Tokens > float64(tb.config.Burst) {
		cl.Tokens = float64(tb.config.Burst)
	}

	// Update the last refill time
	cl.LastRefill = now
}

func (tb *TokenBucket) consumeToken(cl *client) bool {
	if cl.Tokens-1 < 0 {
		cl.Tokens = 0
		return false
	}

	cl.Tokens--

	return true
}
