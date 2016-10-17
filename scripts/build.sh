#!/bin/bash

set -e

OS="darwin linux windows"
ARCH="amd64 386"

rm -Rf bin/
mkdir bin/

for GOOS in $OS; do
    for GOARCH in $ARCH; do
        arch="$GOOS-$GOARCH"
        binary="packer-builder-sakuracloud"
        if [ "$GOOS" = "windows" ]; then
          binary="${binary}.exe"
        fi
        echo "Building $binary $arch"
        GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 \
            govendor build \
                -ldflags "-s -w" \
                -o $binary \
                main.go
        zip -r "bin/packer-builder-sakuracloud_$arch" $binary
        rm -f $binary
    done
done
