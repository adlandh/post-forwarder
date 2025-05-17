package driver

//go:generate go tool oapi-codegen -config ../../../.codegen.yml "../../../api/post-forwarder.yaml"
//go:generate go tool gowrap gen -p github.com/adlandh/post-forwarder/internal/post-forwarder/driver -i ServerInterface -t https://raw.githubusercontent.com/adlandh/gowrap-templates/main/echo-sentry.gotmpl -o open_api_sentry.gen.go -l "" -g
