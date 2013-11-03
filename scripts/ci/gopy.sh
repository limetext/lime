#!/usr/bin/env sh

cd ./3rdparty/libs/gopy
cat lib/cgo.go | python -c "import re,sys; print re.sub(r'(CFLAGS:[^\n]+)(\n+.*?)(LDFLAGS:[^\n]+)', 'CFLAGS: `python3.3-config --cflags`\g<2>LDFLAGS: -L/usr/local/lib `python3.3-config --ldflags`', sys.stdin.read())" > lib/cgo2.go && rm -f lib/cgo.go
cat lib/cgo2.go
go test -i
go test
