package domain

import (
	_ "github.com/adlandh/gowrap-templates/helpers/sentry"
	_ "github.com/oapi-codegen/runtime"
)

//go:generate ../../../gen-wraps.sh
//go:generate go tool mockery
