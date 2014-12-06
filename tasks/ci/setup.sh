#!/usr/bin/env bash

# Just so that our oniguruma.pc is found if
# the user doesn't have an oniguruma.pc.
export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:$PWD/../../../github.com/limetext/rubex

# Colors.
export RED="\e[31m"
export GREEN="\e[32m"
export YELLOW="\e[33m"
export RESET="\e[0m"

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

export -f fold_start
export -f fold_end
