#!/bin/bash

COMMIT_ID=$(git rev-parse HEAD)
BUILD_TS=$(date +'%Y%m%d%H%M%S')
BRANCH_TAG=$(git rev-parse --abbrev-ref HEAD)
BUILD=build-$BUILD_TS-$BRANCH_TAG-$COMMIT_ID

echo "docker building usms..."

docker buildx build \
 --platform linux/amd64 \
 --build-arg BUILD=$BUILD \
 -f build/server/Dockerfile \
 -t usms .