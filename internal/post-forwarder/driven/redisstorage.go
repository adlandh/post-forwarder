package driven

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/adlandh/post-forwarder/internal/post-forwarder/config"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

const ttl = 24 * time.Hour

var _ domain.MessageStorage = (*RedisStorage)(nil)

type RedisStorage struct {
	client *redis.Client
	prefix string
}

func NewRedisStorage(lc fx.Lifecycle, cfg *config.Config) (*RedisStorage, error) {
	opt, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		return nil, fmt.Errorf("error parsing redis url: %w", err)
	}

	r := &RedisStorage{
		client: redis.NewClient(opt),
		prefix: cfg.Redis.Prefix,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := r.client.Set(ctx, r.prefix+"::ping", "pong", 1*time.Millisecond).Err()
			if err != nil {
				return fmt.Errorf("error connecting to redis: %w", err)
			}
			return nil
		},
	})

	return r, nil
}

func (r RedisStorage) Store(ctx context.Context, msg string) (id string, err error) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		err = fmt.Errorf("error generating uuid: %w", err)
		return
	}

	id = newUUID.String()

	err = r.client.Set(ctx, r.prefix+"::"+id, msg, ttl).Err()
	if err != nil {
		err = fmt.Errorf("error storing to redis: %w", err)
	}

	return
}

func (r RedisStorage) Read(ctx context.Context, id string) (msg string, createdAt time.Time, err error) {
	id = r.prefix + "::" + id

	msg, err = r.client.Get(ctx, id).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = domain.ErrorNotFound

			return
		}

		err = fmt.Errorf("error reading from redis: %w", err)

		return
	}

	remainingTTL := r.client.TTL(ctx, id).Val()
	createdAt = time.Now().Add(remainingTTL - ttl)

	return
}
