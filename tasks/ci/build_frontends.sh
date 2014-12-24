#!/usr/bin/env bash

source "$(dirname -- "$0")/setup.sh"

function build {
	pushd "$1"
	go build
	build_result=$?
	echo -ne "${YELLOW}=>${RESET} build $1 - "
	if [ "$build_result" == "0" ]; then
		echo -e "${GREEN}SUCCEEDED${RESET}"
	else
		echo -e "${RED}FAILED${RESET}"
	fi
	popd
}

fold_start "get.frontends" "get frontends"
go get github.com/limetext/lime/frontend/...
fold_end "get.frontends"

ret=0

fold_start "build.termbox" "build termbox frontend"
build "frontend/termbox"
let ret=$ret+$build_result
fold_end "build.termbox"

fold_start "build.html" "build html frontend"
build "frontend/html"
let ret=$ret+$build_result
fold_end "build.html"

fold_start "build.qml" "build qml frontend"
build "frontend/qml"
let ret=$ret+$build_result
fold_end "build.qml"

exit $ret
