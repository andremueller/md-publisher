#!/usr/bin/env bash
set -o errexit
set -o nounset

SCRIPT="$(cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd -P)"
SCRIPT_PATH="$(dirname "$SCRIPT")"
echo $SCRIPT
function die() {
    echo "ERROR $? IN $SCRIPT AT LINE ${BASH_LINENO[0]}"
    exit 1
}
trap die ERR

export GO111MODULE=on

build_opts=()
VERSION="$(git describe --tags --long 2> /dev/null || echo "")"
# returns something like v1.2-3-g177b3eb
[[ -n "$VERSION" ]] && build_opts+=(-ldflags "-X main.version=$VERSION")

# if go is not locally found in path: use container
GO="go"
if ! go version > /dev/null 2>&1 ; then
    echo "go not found in path - using container"
    GO="$SCRIPT_PATH/go"
fi
echo "Using go $GO"
"$GO" get
"$GO" mod tidy
"$GO" build  "${build_opts[@]}" "${@}"

