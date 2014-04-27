#!/usr/bin/env bash

# Just so that our oniguruma.pc is found if
# the user doesn't have an oniguruma.pc.
export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:$PWD/../rubex

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

# Just to fetch all dependencies needed
fold_start "bootstrap"
go get -d -u github.com/limetext/lime/frontend/termbox
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

do_test "backend"
fail1=$build_result

do_test "frontend/termbox"
fail2=$build_result

let ex=$fail1+$fail2
exit $ex
