#!/bin/bash
# Copyright (c) CloudCover
# SPDX-License-Identifier: Apache-2.0

# This script assumes that a tag that is passed in (i.e. $1) has already been pushed
# to github
#
set -ex

VERSION=$1
BASE_DIR=$(pwd)
CLI_NAME=terrarium
CLI_DIR=src/cli/terrarium
BUILD_DATE=$(date "+%Y-%m-%d")
BUILD_DIR=$BASE_DIR/.bin
XGO_GO_VERSION="go-1.20.x"
XGO_TARGETS="darwin/amd64,linux/amd64,darwin/arm64,linux/arm64,windows/amd64"
BIN_TARGETS_IN=("darwin-10.12-amd64" "darwin-10.12-arm64" "linux-amd64" "linux-arm64" "windows-4.0-amd64")
BIN_TARGETS_OUT=("darwin_amd64_v1" "darwin_arm64" "linux_amd64_v1" "linux_arm64" "windows_amd64_v1")

GO_LDFLAGS="-s -w -X github.com/cldcvr/terrarium/src/cli/internal/build.Version=${VERSION} -X github.com/cldcvr/terrarium/src/cli/internal/build.Date=${BUILD_DATE}"

#
# Use xgo to create all the required binaries via cross-compile
#
mkdir -p ${BUILD_DIR}
cd ${BUILD_DIR}
xgo -out ${CLI_NAME}-${VERSION} \
 	-go ${XGO_GO_VERSION}  \
 	--targets=${XGO_TARGETS} \
 	-ldflags="${GO_LDFLAGS}" \
 	github.com/cldcvr/terrarium/${CLI_DIR}

#
# Stage the binaries in the directory structure expected by goreleaser
#
for i in ${!BIN_TARGETS_IN[@]}; do
    ext=""
    if [[ "${BIN_TARGETS_IN[$i]}" == "windows"* ]]; then
      ext=".exe"
    fi
    mkdir ${CLI_NAME}_${BIN_TARGETS_OUT[$i]}
    mv ${CLI_NAME}-${VERSION}-${BIN_TARGETS_IN[$i]}${ext} ${CLI_NAME}_${BIN_TARGETS_OUT[$i]}/${CLI_NAME}${ext}
done

#
# Use our forked goreleaser generate the required artifacts and release the terrarium CLI
#
cd ${BASE_DIR}/${CLI_DIR}
goreleaser release --clean

