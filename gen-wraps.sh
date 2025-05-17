#!/bin/sh

for i in `go tool find-interfaces -path .`
do
  go tool gowrap gen -i $i -t https://raw.githubusercontent.com/adlandh/gowrap-templates/main/sentry.gotmpl -o ./wrappers/${i}WithSentry.gen.go -l "" -g
done
