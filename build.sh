#!/usr/bin/env bash
set -o errexit
set -o nounset

SCRIPT_PATH="$(cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd -P)"
function die() {
    echo "ERROR $? IN ${BASH_SOURCE[0]} AT LINE ${BASH_LINENO[0]}"
    exit 1
}
trap die ERR

export GO111MODULE=on

build_opts=()
VERSION="$(git describe --tags --long 2> /dev/null || echo "")"
# returns something like v1.2-3-g177b3eb
[[ -n "$VERSION" ]] && build_opts+=(-ldflags "-X main.version=$VERSION")
GO="go"
"$GO" clean
"$GO" get -u
"$GO" mod tidy
"$GO" build  "${build_opts[@]}"

echo "-------------- Running unit tests -------------- "
"$GO" clean -testcache ./...
"$GO" test -v ./...
