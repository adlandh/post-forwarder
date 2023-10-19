package driven

import (
	"context"
	"testing"
	"time"

	"github.com/adlandh/post-forwarder/internal/post-forwarder/config"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/fx/fxtest"
)

type RedisStorageTestSuite struct {
	suite.Suite
	storage *RedisStorage
}

func (s *RedisStorageTestSuite) SetupSuite() {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("Ready to accept connections").WithStartupTimeout(3*time.Minute),
			wait.ForListeningPort("6379/tcp").WithStartupTimeout(3*time.Minute),
		),
		Name: "redis",
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	s.Require().NoError(err)

	redisPort, err := container.MappedPort(ctx, "6379")
	s.Require().NoError(err)

	port := redisPort.Port()

	host, err := container.Host(ctx)
	s.Require().NoError(err)

	if host == "" {
		host = "localhost"
	}

	lc := fxtest.NewLifecycle(s.T())

	s.storage, err = NewRedisStorage(
		lc,
		&config.Config{
			Redis: config.RedisConfig{
				URL:    "redis://" + host + ":" + port,
				Prefix: gofakeit.Word(),
			},
		})
	s.Require().NoError(err)

	err = lc.Start(ctx)
	s.Require().NoError(err)
}

func (s *RedisStorageTestSuite) TestStore() {
	msg := gofakeit.SentenceSimple()

	id, err := s.storage.Store(context.Background(), msg)
	s.Require().NoError(err)
	s.Require().NotEmpty(id)
}

func (s *RedisStorageTestSuite) TestRead() {
	msg := gofakeit.SentenceSimple()

	id, err := s.storage.Store(context.Background(), msg)
	s.Require().NoError(err)
	s.Require().NotEmpty(id)
	createdAt := time.Now()

	storedMsg, storedCreatedAt, err := s.storage.Read(context.Background(), id)
	s.Require().NoError(err)
	s.Require().Equal(msg, storedMsg)
	s.Require().Equal(createdAt.Unix(), storedCreatedAt.Unix())
}

func (s *RedisStorageTestSuite) TestNotFound() {
	id := gofakeit.UUID()

	storedMsg, storedCreatedAt, err := s.storage.Read(context.Background(), id)
	s.Require().Error(err)
	s.Require().Equal(domain.ErrorNotFound, err)
	s.Require().Empty(storedMsg)
	s.Require().Empty(storedCreatedAt)
}

func TestRedisStorage(t *testing.T) {
	suite.Run(t, new(RedisStorageTestSuite))
}
