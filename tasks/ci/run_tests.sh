#!/usr/bin/env bash

source "$(dirname -- "$0")/setup.sh"

function run_tests {
	go test "$1" -covermode=count -coverprofile=tmp.cov
	test_result=$?
	# Can't do race tests at the same time as coverage as it'll report
	# lots of false positives
	go test -race "$1"
	let test_result=$test_result+$?
	echo -ne "${YELLOW}=>${RESET} test $1 - "
	if [ "$test_result" == "0" ]; then
		echo -e "${GREEN}SUCCEEDED${RESET}"
	else
		echo -e "${RED}FAILED${RESET}"
	fi
}

function test_all {
	let a=0
	for pkg in $(go list "./$1/..."); do
		run_tests "$pkg"
		let a=$a+$test_result
		if [ "$test_result" == "0" ]; then
			sed 1d tmp.cov >> coverage.cov
		fi
		rm tmp.cov
	done
	test_result=$a
}

fold_start "get.cov" "get coverage tools"
go get golang.org/x/tools/cmd/cover
go get github.com/mattn/goveralls
go get github.com/axw/gocov/gocov
fold_end "get.cov"

fold_start "get.termbox" "get termbox"
go get github.com/limetext/lime/frontend/termbox
fold_end "get.termbox"

echo "mode: count" > coverage.cov

ret=0

fold_start "test.backend" "test backend"
test_all "backend"
let ret=$ret+$test_result
fold_end "test.backend"

fold_start "test.termbox" "test termbox"
test_all "frontend/termbox"
let ret=$ret+$test_result
fold_end "test.termbox"

if [ "$ret" == "0" ] && [ "$TRAVIS_OS_NAME" == "linux" ]; then
	fold_start "coveralls" "post to coveralls"
	"$(go env GOPATH | awk 'BEGIN{FS=":"} {print $1}')/bin/goveralls" -coverprofile=coverage.cov -service=travis-ci
	fold_end "coveralls"
fi

rm coverage.cov

exit $ret
