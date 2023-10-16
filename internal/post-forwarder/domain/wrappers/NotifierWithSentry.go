// Code generated by gowrap. DO NOT EDIT.
// template: https://raw.githubusercontent.com/adlandh/gowrap-templates/main/sentry.gotmpl
// gowrap: http://github.com/hexdigest/gowrap

package wrappers

//go:generate gowrap gen -p github.com/adlandh/post-forwarder/internal/post-forwarder/domain -i Notifier -t https://raw.githubusercontent.com/adlandh/gowrap-templates/main/sentry.gotmpl -o NotifierWithSentry.go -l ""

import (
	"context"

	helpers "github.com/adlandh/gowrap-templates/helpers/sentry"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain"
	"github.com/getsentry/sentry-go"
)

// NotifierWithSentry implements domain.Notifier interface instrumented with opentracing spans
type NotifierWithSentry struct {
	domain.Notifier
	_spanDecorator func(span *sentry.Span, params, results map[string]interface{})
	_instance      string
}

// NewNotifierWithSentry returns NotifierWithSentry
func NewNotifierWithSentry(base domain.Notifier, instance string, spanDecorator ...func(span *sentry.Span, params, results map[string]interface{})) NotifierWithSentry {
	d := NotifierWithSentry{
		Notifier:  base,
		_instance: instance,
	}

	if len(spanDecorator) > 0 && spanDecorator[0] != nil {
		d._spanDecorator = spanDecorator[0]
	} else {
		d._spanDecorator = helpers.SpanDecorator
	}

	return d
}

// Send implements domain.Notifier
func (_d NotifierWithSentry) Send(ctx context.Context, service string, msg string) (err error) {
	span := sentry.StartSpan(ctx, _d._instance+".domain.Notifier.Send", sentry.WithTransactionName("domain.Notifier.Send"))
	ctx = span.Context()

	defer func() {
		_d._spanDecorator(span, map[string]interface{}{
			"ctx":     ctx,
			"service": service,
			"msg":     msg}, map[string]interface{}{
			"err": err})
		span.Finish()
	}()
	return _d.Notifier.Send(ctx, service, msg)
}
