package driven

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/adlandh/post-forwarder/internal/post-forwarder/config"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/ksuid"
	"go.uber.org/fx"
)

const ttl = 24 * time.Hour

var _ domain.MessageStorage = (*RedisStorage)(nil)

const redisKeySeparator = "::"

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
			if err := r.client.Ping(ctx).Err(); err != nil {
				return fmt.Errorf("error connecting to redis: %w", err)
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return r.client.Close()
		},
	})

	return r, nil
}

func (r RedisStorage) Store(ctx context.Context, msg string) (id string, err error) {
	newID, err := ksuid.NewRandom()
	if err != nil {
		err = fmt.Errorf("error generating id: %w", err)
		return
	}

	id = newID.String()

	err = r.client.Set(ctx, r.key(id), msg, ttl).Err()
	if err != nil {
		err = fmt.Errorf("error storing to redis: %w", err)
	}

	return
}

func (r RedisStorage) Read(ctx context.Context, id string) (msg string, createdAt time.Time, err error) {
	key := r.key(id)
	pipe := r.client.Pipeline()
	getCmd := pipe.Get(ctx, key)
	ttlCmd := pipe.TTL(ctx, key)
	_, execErr := pipe.Exec(ctx)

	if execErr != nil && !errors.Is(execErr, redis.Nil) {
		err = fmt.Errorf("error reading from redis: %w", execErr)
		return
	}

	msg, err = getCmd.Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = domain.ErrorNotFound

			return
		}

		err = fmt.Errorf("error reading from redis: %w", err)

		return
	}

	remainingTTL, ttlErr := ttlCmd.Result()
	if ttlErr != nil {
		err = fmt.Errorf("error reading from redis: %w", ttlErr)
		return
	}

	if remainingTTL > 0 {
		now := time.Now()
		createdAt = now.Add(remainingTTL - ttl)
	}

	return
}

func (r RedisStorage) key(id string) string {
	return r.prefix + redisKeySeparator + id
}
