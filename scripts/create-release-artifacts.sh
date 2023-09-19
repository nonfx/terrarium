#!/bin/bash
# Copyright (c) CloudCover
# SPDX-License-Identifier: Apache-2.0

set -ex
VERSION=$1

CLI_NAME=terrarium
BUILD_DATE=$(date "+%Y-%m-%d")
BUILD_DIR=$(pwd)/.bin
XGO_GO_VERSION="go-1.20.x"
XGO_TARGETS="darwin/amd64,linux/amd64,darwin/arm64,linux/arm64,windows/amd64"
GZIP_TARGETS_IN=("darwin-10.12-amd64" "darwin-10.12-arm64" "linux-amd64" "linux-arm64")
GZIP_TARGETS_OUT=("macos-amd64" "macos-arm64" "linux-amd64" "linux-arm64")
ZIP_TARGETS_IN=("windows-4.0-amd64")
ZIP_TARGETS_OUT=("windows-amd64")

GO_LDFLAGS="-s -w -X github.com/cldcvr/terrarium/src/cli/internal/build.Version=${VERSION} -X github.com/cldcvr/terrarium/src/cli/internal/build.Date=${BUILD_DATE}"

mkdir -p ${BUILD_DIR}
cd ${BUILD_DIR}
xgo -out ${CLI_NAME}-${VERSION} \
 	-go ${XGO_GO_VERSION}  \
 	--targets=${XGO_TARGETS} \
 	-ldflags="${GO_LDFLAGS}" \
 	github.com/cldcvr/terrarium/src/cli/terrarium

for i in ${!GZIP_TARGETS_IN[@]}; do
    mv ${CLI_NAME}-${VERSION}-${GZIP_TARGETS_IN[$i]} ${CLI_NAME}
    tar -czvf ${CLI_NAME}-${VERSION}-${GZIP_TARGETS_OUT[$i]}.tar.gz ${CLI_NAME}
done
rm ${CLI_NAME}

for i in ${!ZIP_TARGETS_IN[@]}; do
    mv ${CLI_NAME}-${VERSION}-${ZIP_TARGETS_IN[$i]}.exe ${CLI_NAME}.exe
    zip -r ${CLI_NAME}-${VERSION}-${ZIP_TARGETS_OUT[$i]}.zip ${CLI_NAME}.exe
done
rm ${CLI_NAME}.exe


