// Code generated by gowrap. DO NOT EDIT.
// template: https://raw.githubusercontent.com/adlandh/gowrap-templates/main/sentry.gotmpl
// gowrap: http://github.com/hexdigest/gowrap

package wrappers

//go:generate gowrap gen -p github.com/adlandh/post-forwarder/internal/post-forwarder/domain -i MessageStorage -t https://raw.githubusercontent.com/adlandh/gowrap-templates/main/sentry.gotmpl -o MessageStorageWithSentry.go -l ""

import (
	"context"
	"time"

	helpers "github.com/adlandh/gowrap-templates/helpers/sentry"
	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain"
	"github.com/getsentry/sentry-go"
)

// MessageStorageWithSentry implements domain.MessageStorage interface instrumented with opentracing spans
type MessageStorageWithSentry struct {
	domain.MessageStorage
	_spanDecorator func(span *sentry.Span, params, results map[string]interface{})
	_instance      string
}

// NewMessageStorageWithSentry returns MessageStorageWithSentry
func NewMessageStorageWithSentry(base domain.MessageStorage, instance string, spanDecorator ...func(span *sentry.Span, params, results map[string]interface{})) MessageStorageWithSentry {
	d := MessageStorageWithSentry{
		MessageStorage: base,
		_instance:      instance,
	}

	if len(spanDecorator) > 0 && spanDecorator[0] != nil {
		d._spanDecorator = spanDecorator[0]
	} else {
		d._spanDecorator = helpers.SpanDecorator
	}

	return d
}

// Read implements domain.MessageStorage
func (_d MessageStorageWithSentry) Read(ctx context.Context, id string) (msg string, createdAt time.Time, err error) {
	span := sentry.StartSpan(ctx, _d._instance+".domain.MessageStorage.Read", sentry.WithTransactionName("domain.MessageStorage.Read"))
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
	return _d.MessageStorage.Read(ctx, id)
}

// Store implements domain.MessageStorage
func (_d MessageStorageWithSentry) Store(ctx context.Context, msg string) (id string, err error) {
	span := sentry.StartSpan(ctx, _d._instance+".domain.MessageStorage.Store", sentry.WithTransactionName("domain.MessageStorage.Store"))
	ctx = span.Context()

	defer func() {
		_d._spanDecorator(span, map[string]interface{}{
			"ctx": ctx,
			"msg": msg}, map[string]interface{}{
			"id":  id,
			"err": err})
		span.Finish()
	}()
	return _d.MessageStorage.Store(ctx, msg)
}