#!/usr/bin/env bash

set -euo pipefail

pluginDir=".semrel/$(go env GOOS)_$(go env GOARCH)/changelog-generator-default/0.0.0-dev/"
[[ ! -d "$pluginDir" ]] && {
  echo "creating $pluginDir"
  mkdir -p "$pluginDir"
}

go build -o "$pluginDir/changelog-generator-default" ./cmd/changelog-generator-default
