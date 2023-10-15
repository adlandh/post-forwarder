package main

import (
	"context"
	"time"

	echo_sentry_middleware "github.com/adlandh/echo-sentry-middleware"
	echo_zap_middleware "github.com/adlandh/echo-zap-middleware"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/application"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/config"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain/wrappers"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/driven"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/driver"
	sentry_zapcore "github.com/adlandh/sentry-zapcore"
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func NewApplication(cfg *config.Config, notifier domain.Notifier, logger *zap.Logger) *application.Application {
	if cfg.Sentry.DSN != "" {
		logger = logger.WithOptions(sentry_zapcore.WithSentryOption(sentry_zapcore.WithStackTrace()))
	}

	return application.NewApplication(notifier, logger)
}

func NewSentry(lc fx.Lifecycle, cfg *config.Config) error {
	if cfg.Sentry.DSN == "" {
		return nil
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			return sentry.Init(sentry.ClientOptions{
				Dsn:                cfg.Sentry.DSN,
				EnableTracing:      true,
				TracesSampleRate:   cfg.Sentry.TracesSampleRate,
				ProfilesSampleRate: cfg.Sentry.ProfilesSampleRate,
				MaxErrorDepth:      1,
				Environment:        cfg.Sentry.Environment,
			})
		},
		OnStop: func(ctx context.Context) error {
			sentry.Flush(2 * time.Second)

			return nil
		},
	})

	return nil
}

func NewEcho(lc fx.Lifecycle, server driver.ServerInterface, cfg *config.Config, log *zap.Logger) *echo.Echo {
	e := echo.New()
	e.Use(echo_zap_middleware.MiddlewareWithConfig(log, echo_zap_middleware.ZapConfig{
		AreHeadersDump: true,
		IsBodyDump:     true,
		LimitHTTPBody:  true,
		LimitSize:      1024,
	}))
	e.Use(middleware.Secure())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("1M"))
	e.Use(middleware.RequestID())
	if cfg.Sentry.DSN != "" {
		e.Use(sentryecho.New(sentryecho.Options{
			Repanic: true,
		}))
		e.Use(echo_sentry_middleware.MiddlewareWithConfig(echo_sentry_middleware.SentryConfig{
			AreHeadersDump: true,
			IsBodyDump:     true,
		}))
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			driver.RegisterHandlers(e, server)
			go func() {
				err = e.Start(":" + cfg.Port)
			}()
			return err
		},
		OnStop: func(ctx context.Context) error {
			return e.Shutdown(ctx)
		},
	})

	return e
}

func DecorateServerInterface(cfg *config.Config, srv driver.ServerInterface) driver.ServerInterface {
	if cfg.Sentry.DSN != "" {
		srv = driver.NewServerInterfaceWithSentry(srv, "handlers")
	}

	return srv
}

func DecorateApplicationInterface(cfg *config.Config, app domain.ApplicationInterface) domain.ApplicationInterface {
	if cfg.Sentry.DSN != "" {
		app = wrappers.NewApplicationInterfaceWithSentry(app, "application")
	}

	return app
}

func DecorateNotifier(cfg *config.Config, md domain.Notifier) domain.Notifier {
	if cfg.Sentry.DSN != "" {
		md = wrappers.NewNotifierWithSentry(md, "notifier")
	}

	return md
}

func main() {
	fx.New(CreateService()).Run()
}

func CreateService() fx.Option {
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
				driven.NewNotifiers,
				fx.As(new(domain.Notifier)),
			),
			fx.Annotate(
				NewApplication,
				fx.As(new(domain.ApplicationInterface)),
			),
			fx.Annotate(
				driver.NewHttpServer,
				fx.As(new(driver.ServerInterface)),
			),
		),
		fx.Decorate(
			DecorateServerInterface,
			DecorateApplicationInterface,
			DecorateNotifier,
		),
		fx.Invoke(
			NewSentry,
			NewEcho,
		),
	)
}
