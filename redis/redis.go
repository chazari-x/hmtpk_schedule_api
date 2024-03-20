package redis

import (
	"time"

	"github.com/chazari-x/hmtpk_schedule_api/config"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

func Connect(cfg *config.Redis) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + cfg.Port,
		Password: cfg.Pass,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	return client, client.Ping(ctx).Err()
}
