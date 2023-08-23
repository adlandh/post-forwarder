package main

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestCreateService(t *testing.T) {
	err := fx.ValidateApp(CreateService())
	require.NoError(t, err)
}
