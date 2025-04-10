package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	contextlogger "github.com/adlandh/context-logger"
	sentryExtractor "github.com/adlandh/context-logger/sentry-extractor"
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

func newLogger(cfg *config.Config) (*zap.Logger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("error creating logger: %w", err)
	}

	if cfg.Sentry.DSN != "" {
		logger = logger.WithOptions(sentryZapcore.WithSentryOption(sentryZapcore.WithStackTrace()))
	}

	return logger, nil
}

func newContextLogger(logger *zap.Logger) *contextlogger.ContextLogger {
	return contextlogger.WithContext(logger, contextlogger.WithValueExtractor(domain.RequestID), sentryExtractor.With())
}

func newSentry(lc fx.Lifecycle, cfg *config.Config) error {
	if cfg.Sentry.DSN == "" {
		return nil
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			err := sentry.Init(sentry.ClientOptions{
				Dsn:              cfg.Sentry.DSN,
				EnableTracing:    true,
				TracesSampleRate: cfg.Sentry.TracesSampleRate,
				MaxErrorDepth:    1,
				Environment:      cfg.Sentry.Environment,
			})

			if err != nil {
				return fmt.Errorf("error initializing sentry: %w", err)
			}

			return nil
		},
		OnStop: func(_ context.Context) error {
			sentry.Flush(2 * time.Second)

			return nil
		},
	})

	return nil
}

func newEcho(lc fx.Lifecycle, server driver.ServerInterface, cfg *config.Config, logger *contextlogger.ContextLogger) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Configure middleware with optimized settings
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         4 << 10,
		DisableStackAll:   true,
		DisablePrintStack: true,
	}))
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "SAMEORIGIN",
		HSTSMaxAge:         31536000,
		HSTSPreloadEnabled: true,
	}))
	e.Use(middleware.BodyLimit("1M"))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level:     5,
		MinLength: 256,
	}))
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

	e.Use(echoZapMiddleware.MiddlewareWithContextLogger(logger))

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			driver.RegisterHandlers(e, server)
			go func() {
				if err := e.Start(":" + cfg.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Ctx(ctx).Error("error starting echo server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			if err := e.Shutdown(ctx); err != nil {
				return fmt.Errorf("error shutting down echo server: %w", err)
			}
			return nil
		},
	})

	return e
}

func main() {
	fx.New(createService()).Run()
}

func createService() fx.Option {
	options := fx.Options(
		fx.WithLogger(
			func(log *zap.Logger) fxevent.Logger {
				return &fxevent.ZapLogger{Logger: log}
			},
		),
		fx.Provide(
			config.NewConfig,
			newLogger,
			newContextLogger,
			fx.Annotate(
				driven.NewRedisStorage,
				fx.As(new(domain.MessageStorage)),
			),
			fx.Annotate(
				driven.NewNotifiers,
				fx.As(new(domain.Notifier)),
			),
			fx.Annotate(
				application.NewApplication,
				fx.As(new(domain.ApplicationInterface)),
			),
			fx.Annotate(
				driver.NewHTTPServer,
				fx.As(new(driver.ServerInterface)),
			),
		),
		fx.Invoke(
			newSentry,
			newEcho,
		),
	)

	if os.Getenv("SENTRY_DSN") != "" {
		options = fx.Options(options, fx.Decorate(
			driver.DecorateServerInterfaceWithSentry,
			wrappers.DecorateApplicationInterfaceWithSentry,
			wrappers.DecorateNotifierWithSentry,
			wrappers.DecorateMessageStorageWithSentry,
		))
	}

	return options
}
