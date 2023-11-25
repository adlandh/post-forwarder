package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	echoSentryMiddleware "github.com/adlandh/echo-sentry-middleware"
	echoZapMiddleware "github.com/adlandh/echo-zap-middleware"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/application"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/config"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain/wrappers"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/driven"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/driver"
	sentryZapcore "github.com/adlandh/sentry-zapcore"
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func newApplication(cfg *config.Config, notifier domain.Notifier, logger *zap.Logger, storage domain.MessageStorage) *application.Application {
	if cfg.Sentry.DSN != "" {
		logger = logger.WithOptions(sentryZapcore.WithSentryOption(sentryZapcore.WithStackTrace()))
	}

	return application.NewApplication(notifier, logger, storage)
}

func newSentry(lc fx.Lifecycle, cfg *config.Config) error {
	if cfg.Sentry.DSN == "" {
		return nil
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			err := sentry.Init(sentry.ClientOptions{
				Dsn:                cfg.Sentry.DSN,
				EnableTracing:      true,
				TracesSampleRate:   cfg.Sentry.TracesSampleRate,
				ProfilesSampleRate: cfg.Sentry.ProfilesSampleRate,
				MaxErrorDepth:      1,
				Environment:        cfg.Sentry.Environment,
			})

			if err != nil {
				return fmt.Errorf("error initializing sentry: %w", err)
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			sentry.Flush(2 * time.Second)

			return nil
		},
	})

	return nil
}

func newEcho(lc fx.Lifecycle, server driver.ServerInterface, cfg *config.Config, log *zap.Logger) *echo.Echo {
	e := echo.New()
	e.Use(echoZapMiddleware.Middleware(log))
	e.Use(middleware.Secure())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("1M"))
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		RequestIDHandler: domain.RequestID.Saver,
	}))

	if cfg.Sentry.DSN != "" {
		e.Use(sentryecho.New(sentryecho.Options{
			Repanic: true,
		}))
		e.Use(echoSentryMiddleware.MiddlewareWithConfig(echoSentryMiddleware.SentryConfig{
			AreHeadersDump: true,
			IsBodyDump:     true,
		}))
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			driver.RegisterHandlers(e, server)
			go func() {
				err = e.Start(":" + cfg.Port)
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Error("error starting echo server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			err := e.Shutdown(ctx)
			if err != nil {
				return fmt.Errorf("error shutting down echo server: %w", err)
			}
			return nil
		},
	})

	return e
}

func decorateServerInterface(cfg *config.Config, srv driver.ServerInterface) driver.ServerInterface {
	if cfg.Sentry.DSN != "" {
		srv = driver.NewServerInterfaceWithSentry(srv, "handlers")
	}

	return srv
}

func decorateApplicationInterface(cfg *config.Config, app domain.ApplicationInterface) domain.ApplicationInterface {
	if cfg.Sentry.DSN != "" {
		app = wrappers.NewApplicationInterfaceWithSentry(app, "application")
	}

	return app
}

func decorateNotifier(cfg *config.Config, md domain.Notifier) domain.Notifier {
	if cfg.Sentry.DSN != "" {
		md = wrappers.NewNotifierWithSentry(md, "notifier")
	}

	return md
}

func decorateMessageStorage(cfg *config.Config, ms domain.MessageStorage) domain.MessageStorage {
	if cfg.Sentry.DSN != "" {
		ms = wrappers.NewMessageStorageWithSentry(ms, "redis")
	}

	return ms
}

func main() {
	fx.New(createService()).Run()
}

func createService() fx.Option {
	return fx.Options(
		fx.WithLogger(
			func(log *zap.Logger) fxevent.Logger {
				return &fxevent.ZapLogger{Logger: log}
			},
		),
		fx.Provide(
			config.NewConfig,
			fx.Annotate(
				zap.NewDevelopment,
			),
			fx.Annotate(
				driven.NewRedisStorage,
				fx.As(new(domain.MessageStorage)),
			),
			fx.Annotate(
				driven.NewNotifiers,
				fx.As(new(domain.Notifier)),
			),
			fx.Annotate(
				newApplication,
				fx.As(new(domain.ApplicationInterface)),
			),
			fx.Annotate(
				driver.NewHTTPServer,
				fx.As(new(driver.ServerInterface)),
			),
		),
		fx.Decorate(
			decorateServerInterface,
			decorateApplicationInterface,
			decorateNotifier,
			decorateMessageStorage,
		),
		fx.Invoke(
			newSentry,
			newEcho,
		),
	)
}
