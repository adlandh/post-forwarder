package main

import (
	"context"
	"strconv"
	"time"

	"github.com/adlandh/echo-sentry-middleware"
	"github.com/adlandh/echo-zap-middleware"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/application"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/config"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain/wrappers"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/driven"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/driver"
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	if cfg.SentryDSN == "" {
		return logger, nil
	}

	return logger.WithOptions(zap.Hooks(
		func(entry zapcore.Entry) error {
			if entry.Level == zapcore.ErrorLevel {
				defer sentry.Flush(2 * time.Second)
				localHub := sentry.CurrentHub().Clone()
				localHub.ConfigureScope(func(scope *sentry.Scope) {
					scope.SetTag("File", entry.Caller.File)
					scope.SetTag("Line", strconv.Itoa(entry.Caller.Line))
					scope.SetLevel(sentry.LevelError)
				})
				localHub.CaptureMessage(entry.Message)
			}
			return nil
		},
	)), nil

}

func NewSentry(lc fx.Lifecycle, cfg *config.Config) error {
	if cfg.SentryDSN == "" {
		return nil
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			return sentry.Init(sentry.ClientOptions{
				Dsn:                cfg.SentryDSN,
				EnableTracing:      true,
				TracesSampleRate:   cfg.SentryTracesSampleRate,
				ProfilesSampleRate: cfg.SentryProfilesSampleRate,
				MaxErrorDepth:      1,
				Environment:        cfg.SentryEnvironment,
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
	if cfg.SentryDSN != "" {
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

func DecorateServerInterface(cfg *config.Config, srv driver.ServerInterface, log *zap.Logger) driver.ServerInterface {
	srv = driver.NewServerInterfaceWithZap(srv, log)
	if cfg.SentryDSN != "" {
		srv = driver.NewServerInterfaceWithSentry(srv, "handlers")
	}

	return srv
}

func DecorateApplicationInterface(cfg *config.Config, app domain.ApplicationInterface, log *zap.Logger) domain.ApplicationInterface {
	app = wrappers.NewApplicationInterfaceWithZap(app, log)
	if cfg.SentryDSN != "" {
		app = wrappers.NewApplicationInterfaceWithSentry(app, "application")
	}

	return app
}

func DecorateMessageDestination(cfg *config.Config, md domain.MessageDestination, log *zap.Logger) domain.MessageDestination {
	md = wrappers.NewMessageDestinationWithZap(md, log)
	if cfg.SentryDSN != "" {
		md = wrappers.NewMessageDestinationWithSentry(md, "bot")
	}

	return md
}

func CreateService() fx.Option {
	return fx.Options(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(
			config.NewConfig,
			NewLogger,
			fx.Annotate(
				driven.NewTelegramMessageSender,
				fx.As(new(domain.MessageDestination)),
			),
			fx.Annotate(
				application.NewApplication,
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
			DecorateMessageDestination,
		),
		fx.Invoke(
			NewSentry,
			NewEcho,
		),
	)
}

func main() {
	fx.New(CreateService()).Run()
}
