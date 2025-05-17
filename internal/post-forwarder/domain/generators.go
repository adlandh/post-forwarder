package domain

//go:generate go tool gowrap gen -i ApplicationInterface -t https://raw.githubusercontent.com/adlandh/gowrap-templates/main/sentry.gotmpl -o ./wrappers/ApplicationInterfaceWithSentry.go -l "" -g -v InstanceName=application
//go:generate go tool gowrap gen -i MessageStorage -t https://raw.githubusercontent.com/adlandh/gowrap-templates/main/sentry.gotmpl -o ./wrappers/MessageStorageWithSentry.go -l "" -g -v InstanceName=notifier
//go:generate go tool gowrap gen -i Notifier -t https://raw.githubusercontent.com/adlandh/gowrap-templates/main/sentry.gotmpl -o ./wrappers/NotifierWithSentry.go -l "" -g -v InstanceName=redis
//go:generate go tool mockery
