#!/usr/bin/env bash

set -euo pipefail

pluginDir=".semrel/$(go env GOOS)_$(go env GOARCH)/hooks-plugin-registry-update/0.0.0-dev/"
[[ ! -d "$pluginDir" ]] && {
  echo "creating $pluginDir"
  mkdir -p "$pluginDir"
}

go build -o "$pluginDir/hooks-plugin-registry-update" ./
