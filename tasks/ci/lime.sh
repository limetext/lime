#!/usr/bin/env bash

# Just so that our oniguruma.pc is found if
# the user doesn't have an oniguruma.pc.
export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:$PWD/../../../github.com/limetext/rubex

# Colors.
RED="\e[31m"
GREEN="\e[32m"
YELLOW="\e[33m"
RESET="\e[0m"

function fold_start {
    if [ "$TRAVIS" == "true" ]; then
        echo -en "travis_fold:start:$1\r"
        echo "\$ $2"
    fi
}

function fold_end {
    if [ "$TRAVIS" == "true" ]; then
        echo -en "travis_fold:end:$1\r"
    fi
}

function run_tests {
	go test "$1" -covermode=count -coverprofile=tmp.cov
	build_result=$?
	# Can't do race tests at the same time as coverage as it'll report
	# lots of false positives then..
	go test -race "$1"
	let build_result=$build_result+$?
	echo -ne "${YELLOW}=>${RESET} test $1 - "
	if [ "$build_result" == "0" ]; then
	    echo -e "${GREEN}SUCCEEDED${RESET}"
	else
	    echo -e "${RED}FAILED${RESET}"
	fi
}

function test_all {
	let a=0
	for pkg in $(go list "./$1/..."); do
		run_tests "$pkg"
		let a=$a+$build_result
		if [ "$build_result" == "0" ]; then
			sed 1d tmp.cov >> coverage.cov
		fi
	done
	build_result=$a
}

fold_start "get_cov" "get coverage tools"
go get code.google.com/p/go.tools/cmd/cover
go get github.com/mattn/goveralls
go get github.com/axw/gocov/gocov
fold_end "get_cov"

fold_start "get_termbox" "get termbox"
go get github.com/limetext/lime/frontend/termbox
fold_end "get_termbox"

echo "mode: count" > coverage.cov

ret=0

fold_start "test_backend" "test backend"
test_all "backend"
let ret=ret+$build_result
fold_end "test_backend"

fold_start "test_termbox" "test termbox"
test_all "frontend/termbox"
let ret=ret+$build_result
fold_end "test_termbox"

if [ "$ret" == "0" ]; then
	fold_start "coveralls" "post to coveralls"
	"$(go env GOPATH | awk 'BEGIN{FS=":"} {print $1}')/bin/goveralls" -coverprofile=coverage.cov -service=travis-ci
	fold_end "coveralls"
fi

exit $ret
