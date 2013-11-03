#!/usr/bin/env sh

cd ./3rdparty/libs/gopy
cat lib/cgo.go | python -c "import re,sys; print re.sub(r'(CFLAGS:[^\n]+)(\n+.*?)(LDFLAGS:[^\n]+)', 'CFLAGS: `python3.3-config --cflags`\g<2>LDFLAGS: `python3.3-config --libs`', sys.stdin.read())" > lib/cgo2.go && rm -f lib/cgo.go
go install
go test
