package floodcontrol

import (
	"context"
)

// FloodControl интерфейс для реализации правил флуд контроля.
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}

type Config struct {
	RPS   float64 // Maximum requests per second
	Burst int     // Maximum burst
}
