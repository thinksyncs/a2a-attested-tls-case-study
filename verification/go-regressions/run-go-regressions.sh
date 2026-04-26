#!/usr/bin/env sh
set -eu

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
repo_root=$(CDPATH= cd -- "$script_dir/../.." && pwd)

if [ "${COCOS_SOURCE:-}" ]; then
  cocos_source=$COCOS_SOURCE
elif [ -f "$repo_root/../go.mod" ] && grep -q "module github.com/ultravioletrs/cocos" "$repo_root/../go.mod"; then
  cocos_source=$(CDPATH= cd -- "$repo_root/.." && pwd)
else
  echo "COCOS_SOURCE is required." >&2
  echo "Example: COCOS_SOURCE=/path/to/cocos $0" >&2
  exit 2
fi

if [ ! -f "$cocos_source/go.mod" ] || ! grep -q "module github.com/ultravioletrs/cocos" "$cocos_source/go.mod"; then
  echo "COCOS_SOURCE does not look like a Cocos checkout: $cocos_source" >&2
  exit 2
fi

workdir=$(mktemp -d "${TMPDIR:-/tmp}/a2a-go-regressions.XXXXXX")
trap 'rm -rf "$workdir"' EXIT

cp "$script_dir/current/authenticator_regression_test.go" "$workdir/authenticator_regression_test.go"

cat > "$workdir/go.mod" <<EOF
module a2a-attested-tls-go-regressions

go 1.26.0

require github.com/ultravioletrs/cocos v0.0.0

replace github.com/ultravioletrs/cocos => $cocos_source
EOF

export GOCACHE=${GOCACHE:-$workdir/gocache}

(cd "$workdir" && go test -run 'Test' -count=1)
