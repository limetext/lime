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
        echo "\$ $1"
    fi
}

function fold_end {
    if [ "$TRAVIS" == "true" ]; then
        echo -en "travis_fold:end:$1\r"
    fi
}

function do_test2 {
	go test "$1" -covermode=count -coverprofile=tmp.cov
	# Can't do race tests at the same time as coverage as it'll report
	# lots of false positives then..
	go test -race "$1"
	build_result=$?
	echo -ne "${YELLOW}=>${RESET} test $1 - "
	if [ "$build_result" == "0" ]; then
	    echo -e "${GREEN}SUCCEEDED${RESET}"
	else
	    echo -e "${RED}FAILED${RESET}"
	fi
}

function do_test {
	let a=0
	for pkg in $(go list "./$1/..."); do
		do_test2 "$pkg"
		let a=$a+$build_result
		cat tmp.cov | sed 1d >> coverage.cov
	done
	build_result=$a
}

fold_start "coverage tools"
go get code.google.com/p/go.tools/cmd/cover
go get github.com/mattn/goveralls
go get github.com/axw/gocov/gocov
fold_end "coverage tools"

# Just to fetch all dependencies needed
# Do *not* add "-u", or pull requests will not actually be built,
# instead it'll checkout master
fold_start "bootstrap"
go get -d github.com/limetext/lime/frontend/termbox
fold_end "bootstrap"

fold_start "Gen Python"
go run tasks/build/gen_python_api.go
fold_end "Gen Python"

fold_start "Add License"
go run tasks/build/fix.go
fold_end "Add License"

# Actual build now that the python api has been refreshed
fold_start "termbox"
go get github.com/limetext/lime/frontend/termbox
fold_end "termbox"

# Installing qml dependencies fails atm.
# fold_start "qml"
# go get github.com/limetext/lime/frontend/qml
# fold_end "qml"

echo "mode: count" > coverage.cov

do_test "backend"
fail1=$build_result

do_test "frontend/termbox"
fail2=$build_result

"$(go env GOPATH | awk 'BEGIN{FS=":"} {print $1}')/bin/goveralls" -coverprofile=coverage.cov -service=travis-ci

let ex=$fail1+$fail2
exit $ex
