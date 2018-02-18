#!/usr/bin/env bash

set -e -u -x

readonly NAME=netatmo-exporter
readonly PACKAGE=github.com/xperimental/netatmo-exporter
readonly BUILDS=(linux/amd64 linux/arm darwin/amd64 windows/amd64)

mkdir -p gopath/src/github.com/xperimental/
ln -s "$PWD/source" "gopath/src/${PACKAGE}"
export GOPATH=$PWD/gopath

for build in "${BUILDS[@]}"; do
  echo "Build: $build"
  IFS="/" read -a buildinfo <<< "$build"
  GOOS=${buildinfo[0]} GOARCH=${buildinfo[1]} go build -v -ldflags="-s -w" -o "binaries/${NAME}_${buildinfo[0]}_${buildinfo[1]}" "${PACKAGE}"
done
