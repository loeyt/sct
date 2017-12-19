#!/bin/bash

set -x

mkdir -p $GOPATH/src/loe.yt
cp -pr sct $GOPATH/src/loe.yt

go install -v loe.yt/sct

if [ "$GOOS" = "windows" ]; then GOEXE=.exe; fi
mv $GOPATH/bin/${GOOS}${GOOS:+_}${GOARCH}/sct${GOEXE} bin/sct-${GOOS:-linux}-${GOARCH:-amd64}${GOARM:+v}${GOARM}${GOEXE}
