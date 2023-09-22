// Code generated by gowrap. DO NOT EDIT.
// template: https://raw.githubusercontent.com/adlandh/gowrap-templates/main/zap.gotmpl
// gowrap: http://github.com/hexdigest/gowrap

package driver

//go:generate gowrap gen -p github.com/adlandh/post-forwarder/internal/post-forwarder/driver -i ServerInterface -t https://raw.githubusercontent.com/adlandh/gowrap-templates/main/zap.gotmpl -o open_api_zap_gen.go -l ""

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// ServerInterfaceWithZap implements ServerInterface that is instrumented with zap logger
type ServerInterfaceWithZap struct {
	_base ServerInterface
	_log  *zap.Logger
}

// NewServerInterfaceWithZap instruments an implementation of the ServerInterface with simple logging
func NewServerInterfaceWithZap(base ServerInterface, log *zap.Logger) ServerInterfaceWithZap {
	return ServerInterfaceWithZap{
		_base: base,
		_log:  log,
	}
}

// GetWebhook implements ServerInterface
func (_d ServerInterfaceWithZap) GetWebhook(ctx echo.Context, token string, service string) (err error) {
	_d._log.Debug("ServerInterfaceWithZap: calling GetWebhook", zap.Any("params", map[string]interface{}{
		"ctx":     ctx,
		"token":   token,
		"service": service}))
	defer func() {
		if err != nil {
			_d._log.Warn("ServerInterfaceWithZap: method GetWebhook returned an error", zap.Error(err), zap.Any("result", map[string]interface{}{
				"err": err}))
		} else {
			_d._log.Debug("ServerInterfaceWithZap: method GetWebhook finished", zap.Any("result", map[string]interface{}{
				"err": err}))
		}
	}()
	return _d._base.GetWebhook(ctx, token, service)
}

// HealthCheck implements ServerInterface
func (_d ServerInterfaceWithZap) HealthCheck(ctx echo.Context) (err error) {
	_d._log.Debug("ServerInterfaceWithZap: calling HealthCheck", zap.Any("params", map[string]interface{}{
		"ctx": ctx}))
	defer func() {
		if err != nil {
			_d._log.Warn("ServerInterfaceWithZap: method HealthCheck returned an error", zap.Error(err), zap.Any("result", map[string]interface{}{
				"err": err}))
		} else {
			_d._log.Debug("ServerInterfaceWithZap: method HealthCheck finished", zap.Any("result", map[string]interface{}{
				"err": err}))
		}
	}()
	return _d._base.HealthCheck(ctx)
}

// PostWebhook implements ServerInterface
func (_d ServerInterfaceWithZap) PostWebhook(ctx echo.Context, token string, service string) (err error) {
	_d._log.Debug("ServerInterfaceWithZap: calling PostWebhook", zap.Any("params", map[string]interface{}{
		"ctx":     ctx,
		"token":   token,
		"service": service}))
	defer func() {
		if err != nil {
			_d._log.Warn("ServerInterfaceWithZap: method PostWebhook returned an error", zap.Error(err), zap.Any("result", map[string]interface{}{
				"err": err}))
		} else {
			_d._log.Debug("ServerInterfaceWithZap: method PostWebhook finished", zap.Any("result", map[string]interface{}{
				"err": err}))
		}
	}()
	return _d._base.PostWebhook(ctx, token, service)
}
