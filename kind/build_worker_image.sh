#/bin/bash

set -Eeuo pipefail
[ ! -z "${DEBUG+x}" ] && set -x

HERE=$(dirname $(readlink -f "$0"))

## Parameters
DOCKERFILE="${HERE}/kind-worker.dockerfile"
IMG=${IMG:-enc-kind-worker}
KIND_IMG="${KIND_IMG:-kindest/node}"
VERSION=${VERSION:-1.22.2}
MIRROR=
# PUSH


FULL_IMG="${IMG}:v${VERSION}"

if ! docker pull "${MIRROR}${FULL_IMG}"; then
    DOCKER_BUILDKIT=1 docker build \
        --build-arg version="$VERSION" \
        --build-arg base="$KIND_IMG" \
        -t "$FULL_IMG" - < "$DOCKERFILE"
    if [[ ! -z "${PUSH+x}" && ! -z "$MIRROR" ]]; then
        docker tag "$FULL_IMG" "${MIRROR}${FULL_IMG}"
        docker push "${MIRROR}${FULL_IMG}"
    fi
fi
