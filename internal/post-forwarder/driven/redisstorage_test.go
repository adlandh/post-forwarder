package driven

import (
	"context"
	"testing"
	"time"

	"github.com/adlandh/post-forwarder/internal/post-forwarder/config"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"go.uber.org/fx/fxtest"
)

type RedisStorageTestSuite struct {
	suite.Suite
	storage        *RedisStorage
	redisContainer *redis.RedisContainer
	lifecycle      *fxtest.Lifecycle
}

const (
	redisTestTimeout = 30 * time.Second
	redisTTLLeeway   = 2 * time.Second
)

func (s *RedisStorageTestSuite) newTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), redisTestTimeout)
}

func (s *RedisStorageTestSuite) SetupSuite() {
	ctx, cancel := s.newTimeoutContext()
	defer cancel()
	redisContainer, err := redis.Run(ctx, "redis:7-alpine")
	s.Require().NoError(err)
	s.redisContainer = redisContainer

	connStr, err := redisContainer.ConnectionString(ctx)
	s.Require().NoError(err)

	lc := fxtest.NewLifecycle(s.T())
	s.lifecycle = lc

	s.storage, err = NewRedisStorage(
		lc,
		&config.Config{
			Redis: config.RedisConfig{
				URL:    connStr,
				Prefix: gofakeit.Word(),
			},
		})
	s.Require().NoError(err)

	err = lc.Start(ctx)
	s.Require().NoError(err)
}

func (s *RedisStorageTestSuite) TearDownSuite() {
	ctx, cancel := s.newTimeoutContext()
	defer cancel()
	if s.lifecycle != nil {
		s.Require().NoError(s.lifecycle.Stop(ctx))
	}
	if s.redisContainer != nil {
		s.Require().NoError(s.redisContainer.Terminate(ctx))
	}
}

func (s *RedisStorageTestSuite) TestStore() {
	msg := gofakeit.Sentence()

	ctx, cancel := s.newTimeoutContext()
	defer cancel()

	id, err := s.storage.Store(ctx, msg)
	s.Require().NoError(err)
	s.Require().NotEmpty(id)

	storedMsg, storedCreatedAt, err := s.storage.Read(ctx, id)
	s.Require().NoError(err)
	s.Require().Equal(msg, storedMsg)
	s.Require().False(storedCreatedAt.IsZero())
}

func (s *RedisStorageTestSuite) TestRead() {
	msg := gofakeit.Sentence()
	beforeStore := time.Now()

	ctx, cancel := s.newTimeoutContext()
	defer cancel()

	id, err := s.storage.Store(ctx, msg)
	s.Require().NoError(err)
	s.Require().NotEmpty(id)
	afterStore := time.Now()

	storedMsg, storedCreatedAt, err := s.storage.Read(ctx, id)
	afterRead := time.Now()
	s.Require().NoError(err)
	s.Require().Equal(msg, storedMsg)
	s.Require().False(storedCreatedAt.IsZero())
	s.Require().False(storedCreatedAt.Before(beforeStore.Add(-redisTTLLeeway)))
	s.Require().False(storedCreatedAt.After(afterStore.Add(redisTTLLeeway)))
	s.Require().False(storedCreatedAt.After(afterRead.Add(redisTTLLeeway)))
}

func (s *RedisStorageTestSuite) TestNotFound() {
	id := gofakeit.UUID()

	ctx, cancel := s.newTimeoutContext()
	defer cancel()

	storedMsg, storedCreatedAt, err := s.storage.Read(ctx, id)
	s.Require().Error(err)
	s.Require().ErrorIs(err, domain.ErrorNotFound)
	s.Require().Empty(storedMsg)
	s.Require().Empty(storedCreatedAt)
}

func (s *RedisStorageTestSuite) TestStoreUsesPrefixAndTTL() {
	msg := gofakeit.Sentence()

	ctx, cancel := s.newTimeoutContext()
	defer cancel()

	id, err := s.storage.Store(ctx, msg)
	s.Require().NoError(err)
	s.Require().NotEmpty(id)

	key := s.storage.prefix + "::" + id
	ttlValue, err := s.storage.client.TTL(ctx, key).Result()
	s.Require().NoError(err)
	s.Require().Greater(ttlValue, time.Duration(0))
	s.Require().LessOrEqual(ttlValue, ttl)
	s.Require().GreaterOrEqual(ttlValue, ttl-redisTTLLeeway)
}

func (s *RedisStorageTestSuite) TestStoreGeneratesUniqueIDs() {
	ctx, cancel := s.newTimeoutContext()
	defer cancel()

	firstID, err := s.storage.Store(ctx, gofakeit.Sentence())
	s.Require().NoError(err)
	s.Require().NotEmpty(firstID)

	secondID, err := s.storage.Store(ctx, gofakeit.Sentence())
	s.Require().NoError(err)
	s.Require().NotEmpty(secondID)

	s.Require().NotEqual(firstID, secondID)
}

func TestRedisStorage(t *testing.T) {
	suite.Run(t, new(RedisStorageTestSuite))
}
