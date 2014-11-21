#!/usr/bin/env sh

fmt="$(find . ! \( -path './packages' -prune \) -type f -name '*.go' -print0 | xargs -0 gofmt -l )"

if [ -n "$fmt" ]; then
    echo "Unformatted Go source code:"
    echo "$fmt"
    exit 1
fi
