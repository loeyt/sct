#!/bin/bash

set -x

mkdir -p $GOPATH/src/loe.yt
cp -pr sct $GOPATH/src/loe.yt

go install -v loe.yt/sct

mv $GOPATH/bin/${GOOS}_${GOARCH}/sct bin/sct-${GOOS:-linux}-${GOARCH:-amd64}${GOARM:+v}${GOARM}
