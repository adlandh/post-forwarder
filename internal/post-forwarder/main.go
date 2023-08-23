package main

import (
	"context"

	"github.com/adlandh/post-forwarder/internal/post-forwarder/application"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/config"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain/wrappers"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/driven"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/driver"
	"github.com/labstack/echo-contrib/echoprometheus"

	"github.com/adlandh/echo-zap-middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func NewEcho(lc fx.Lifecycle, server driver.ServerInterface, cfg *config.Config, log *zap.Logger) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Secure())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("1M"))
	e.Use(echo_zap_middleware.Middleware(log))
	e.Use(middleware.RequestID())
	e.Use(echoprometheus.NewMiddleware("app"))
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			driver.RegisterHandlers(e, server)
			e.GET("/metrics", echoprometheus.NewHandler())
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

func CreateService() fx.Option {
	return fx.Options(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(
			zap.NewDevelopment,
			config.NewConfig,
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
			fx.Annotate(
				wrappers.NewApplicationInterfaceWithZap,
				fx.As(new(domain.ApplicationInterface)),
			),
			fx.Annotate(
				wrappers.NewMessageDestinationWithZap,
				fx.As(new(domain.MessageDestination)),
			),
		),
		fx.Invoke(
			NewEcho,
		),
	)
}

func main() {
	fx.New(CreateService()).Run()
}
