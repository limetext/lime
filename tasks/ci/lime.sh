#!/usr/bin/env bash

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

function do_test {
	go test ./$1/...
	build_result=$?
	echo -ne "${YELLOW}=>${RESET} test $1 - "
	if [ "$build_result" == "0" ]; then
	    echo -e "${GREEN}SUCCEEDED${RESET}"
	else
	    echo -e "${RED}FAILED${RESET}"
	fi
}

fold_start "termbox"
go get github.com/limetext/lime/frontend/termbox
fold_end "termbox"

# Installing qml dependencies fails atm.
# fold_start "qml"
# go get github.com/limetext/lime/frontend/qml
# fold_end "qml"

do_test "backend"
fail1=$build_result

do_test "frontend/termbox"
fail2=$build_result

let ex=$fail1+$fail2
exit $ex
