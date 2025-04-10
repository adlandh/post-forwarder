// Code generated by gowrap. DO NOT EDIT.
// template: https://raw.githubusercontent.com/adlandh/gowrap-templates/main/echo-sentry.gotmpl
// gowrap: http://github.com/hexdigest/gowrap

package driver

import (
	"github.com/labstack/echo/v4"

	helpers "github.com/adlandh/gowrap-templates/helpers/sentry"
	"github.com/getsentry/sentry-go"
)

// ServerInterfaceWithSentry implements ServerInterface interface instrumented with opentracing spans
type ServerInterfaceWithSentry struct {
	ServerInterface
	_spanDecorator func(span *sentry.Span, params, results map[string]interface{})
	_instance      string
}

// NewServerInterfaceWithSentry returns ServerInterfaceWithSentry
func NewServerInterfaceWithSentry(base ServerInterface, instance string, spanDecorator ...func(span *sentry.Span, params, results map[string]interface{})) ServerInterfaceWithSentry {
	if instance == "" {
		instance = "handlers"
	}

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

// DecorateServerInterfaceWithSentry returns ServerInterface with tracing decorators. Useful for uber fx
func DecorateServerInterfaceWithSentry(base ServerInterface) ServerInterface {
	return NewServerInterfaceWithSentry(base, "")
}

// GetWebhook implements ServerInterface
func (_d ServerInterfaceWithSentry) GetWebhook(ctx echo.Context, token string, service string) (err error) {
	request := ctx.Request()
	savedCtx := request.Context()
	span := sentry.StartSpan(savedCtx, _d._instance+".ServerInterface.GetWebhook", sentry.WithTransactionName("ServerInterface.GetWebhook"))
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
	span := sentry.StartSpan(savedCtx, _d._instance+".ServerInterface.HealthCheck", sentry.WithTransactionName("ServerInterface.HealthCheck"))
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
	span := sentry.StartSpan(savedCtx, _d._instance+".ServerInterface.PostWebhook", sentry.WithTransactionName("ServerInterface.PostWebhook"))
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

// ShowMessage implements ServerInterface
func (_d ServerInterfaceWithSentry) ShowMessage(ctx echo.Context, id string) (err error) {
	request := ctx.Request()
	savedCtx := request.Context()
	span := sentry.StartSpan(savedCtx, _d._instance+".ServerInterface.ShowMessage", sentry.WithTransactionName("ServerInterface.ShowMessage"))
	ctxNew := span.Context()

	defer func() {
		_d._spanDecorator(span, map[string]interface{}{
			"ctx": ctx,
			"id":  id}, map[string]interface{}{
			"err": err})
		span.Finish()
	}()
	ctx.SetRequest(request.WithContext(ctxNew))
	return _d.ServerInterface.ShowMessage(ctx, id)
}
