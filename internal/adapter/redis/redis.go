package redis

import (
	"time"

	"github.com/LeseyS/MKK_BASIS_GO_TEST/pkg/redis"
)

const (
	taskListTTL = 5 * time.Minute
)

type Redis struct {
	redis *redis.Client
}

func New(client *redis.Client) *Redis {
	return &Redis{
		redis: client,
	}
}
