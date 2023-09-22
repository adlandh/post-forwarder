// Code generated by gowrap. DO NOT EDIT.
// template: https://raw.githubusercontent.com/adlandh/gowrap-templates/main/echo-sentry.gotmpl
// gowrap: http://github.com/hexdigest/gowrap

package driver

//go:generate gowrap gen -p github.com/adlandh/post-forwarder/internal/post-forwarder/driver -i ServerInterface -t https://raw.githubusercontent.com/adlandh/gowrap-templates/main/echo-sentry.gotmpl -o open_api_sentry_gen.go -l ""

import (
	helpers "github.com/adlandh/gowrap-templates/helpers/sentry"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
)

// ServerInterfaceWithSentry implements ServerInterface interface instrumented with opentracing spans
type ServerInterfaceWithSentry struct {
	ServerInterface
	_instance      string
	_spanDecorator func(span *sentry.Span, params, results map[string]interface{})
}

// NewServerInterfaceWithSentry returns ServerInterfaceWithSentry
func NewServerInterfaceWithSentry(base ServerInterface, instance string, spanDecorator ...func(span *sentry.Span, params, results map[string]interface{})) ServerInterfaceWithSentry {
	d := ServerInterfaceWithSentry{
		ServerInterface: base,
		_instance:       instance,
	}

	if len(spanDecorator) > 0 && spanDecorator[0] != nil {
		d._spanDecorator = spanDecorator[0]
	} else {
		d._spanDecorator = helpers.SpanDecorator
	}

	return d
}

// GetWebhook implements ServerInterface
func (_d ServerInterfaceWithSentry) GetWebhook(ctx echo.Context, token string, service string) (err error) {
	request := ctx.Request()
	savedCtx := request.Context()
	span := sentry.StartSpan(savedCtx, _d._instance+".ServerInterface.GetWebhook", sentry.TransactionName("ServerInterface.GetWebhook"))
	ctxNew := span.Context()

	defer func() {
		_d._spanDecorator(span, map[string]interface{}{
			"ctx":     ctx,
			"token":   token,
			"service": service}, map[string]interface{}{
			"err": err})
		span.Finish()
	}()
	ctx.SetRequest(request.WithContext(ctxNew))
	return _d.ServerInterface.GetWebhook(ctx, token, service)
}

// HealthCheck implements ServerInterface
func (_d ServerInterfaceWithSentry) HealthCheck(ctx echo.Context) (err error) {
	request := ctx.Request()
	savedCtx := request.Context()
	span := sentry.StartSpan(savedCtx, _d._instance+".ServerInterface.HealthCheck", sentry.TransactionName("ServerInterface.HealthCheck"))
	ctxNew := span.Context()

	defer func() {
		_d._spanDecorator(span, map[string]interface{}{
			"ctx": ctx}, map[string]interface{}{
			"err": err})
		span.Finish()
	}()
	ctx.SetRequest(request.WithContext(ctxNew))
	return _d.ServerInterface.HealthCheck(ctx)
}

// PostWebhook implements ServerInterface
func (_d ServerInterfaceWithSentry) PostWebhook(ctx echo.Context, token string, service string) (err error) {
	request := ctx.Request()
	savedCtx := request.Context()
	span := sentry.StartSpan(savedCtx, _d._instance+".ServerInterface.PostWebhook", sentry.TransactionName("ServerInterface.PostWebhook"))
	ctxNew := span.Context()

	defer func() {
		_d._spanDecorator(span, map[string]interface{}{
			"ctx":     ctx,
			"token":   token,
			"service": service}, map[string]interface{}{
			"err": err})
		span.Finish()
	}()
	ctx.SetRequest(request.WithContext(ctxNew))
	return _d.ServerInterface.PostWebhook(ctx, token, service)
}
