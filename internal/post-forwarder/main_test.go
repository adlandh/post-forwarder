package main

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestCreateService(t *testing.T) {
	err := fx.ValidateApp(createService())
	require.NoError(t, err)
}
