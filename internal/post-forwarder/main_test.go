package main

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestCreateService(t *testing.T) {
	t.Setenv("SENTRY_DSN", "test")
	err := fx.ValidateApp(createService())
	require.NoError(t, err)
}
