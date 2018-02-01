#!/usr/bin/env bash

# Execute a command inside the running `wptd-dev-instance`.

DOCKER_DIR=$(dirname "$0")
source "${DOCKER_DIR}/../logging.sh"
source "${DOCKER_DIR}/../path.sh"
WPTD_PATH=${WPTD_PATH:-$(absdir "${DOCKER_DIR}/../..")}

docker exec -t -u $(id -u $USER):$(id -g $USER) wptd-dev-instance "$@"
