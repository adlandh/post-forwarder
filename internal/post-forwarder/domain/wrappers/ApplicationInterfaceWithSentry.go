// Code generated by gowrap. DO NOT EDIT.
// template: https://raw.githubusercontent.com/adlandh/gowrap-templates/main/sentry.gotmpl
// gowrap: http://github.com/hexdigest/gowrap

package wrappers

import (
	"context"
	"time"

	_sourceDomain "github.com/adlandh/post-forwarder/internal/post-forwarder/domain"

	helpers "github.com/adlandh/gowrap-templates/helpers/sentry"
	"github.com/getsentry/sentry-go"
)

// ApplicationInterfaceWithSentry implements _sourceDomain.ApplicationInterface interface instrumented with opentracing spans
type ApplicationInterfaceWithSentry struct {
	_sourceDomain.ApplicationInterface
	_spanDecorator func(span *sentry.Span, params, results map[string]interface{})
	_instance      string
}

// NewApplicationInterfaceWithSentry returns ApplicationInterfaceWithSentry
func NewApplicationInterfaceWithSentry(base _sourceDomain.ApplicationInterface, instance string, spanDecorator ...func(span *sentry.Span, params, results map[string]interface{})) ApplicationInterfaceWithSentry {
	d := ApplicationInterfaceWithSentry{
		ApplicationInterface: base,
		_instance:            instance,
	}

	if len(spanDecorator) > 0 && spanDecorator[0] != nil {
		d._spanDecorator = spanDecorator[0]
	} else {
		d._spanDecorator = helpers.SpanDecorator
	}

	return d
}

// GetMessage implements _sourceDomain.ApplicationInterface
func (_d ApplicationInterfaceWithSentry) GetMessage(ctx context.Context, id string) (msg string, createdAt time.Time, err error) {
	span := sentry.StartSpan(ctx, _d._instance+"._sourceDomain.ApplicationInterface.GetMessage", sentry.WithTransactionName("_sourceDomain.ApplicationInterface.GetMessage"))
	ctx = span.Context()

	defer func() {
		_d._spanDecorator(span, map[string]interface{}{
			"ctx": ctx,
			"id":  id}, map[string]interface{}{
			"msg":       msg,
			"createdAt": createdAt,
			"err":       err})
		span.Finish()
	}()
	return _d.ApplicationInterface.GetMessage(ctx, id)
}

// ProcessRequest implements _sourceDomain.ApplicationInterface
func (_d ApplicationInterfaceWithSentry) ProcessRequest(ctx context.Context, url string, service string, msg string) (err error) {
	span := sentry.StartSpan(ctx, _d._instance+"._sourceDomain.ApplicationInterface.ProcessRequest", sentry.WithTransactionName("_sourceDomain.ApplicationInterface.ProcessRequest"))
	ctx = span.Context()

	defer func() {
		_d._spanDecorator(span, map[string]interface{}{
			"ctx":     ctx,
			"url":     url,
			"service": service,
			"msg":     msg}, map[string]interface{}{
			"err": err})
		span.Finish()
	}()
	return _d.ApplicationInterface.ProcessRequest(ctx, url, service, msg)
}
