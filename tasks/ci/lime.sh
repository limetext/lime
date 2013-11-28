#!/usr/bin/env bash

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

cd build

fold_start "cmake"
cmake ..
fold_end "cmake"

fold_start "make"
make
build_result=$?
fold_end "make"

# Print colored build result.
RED="\e[31m"
GREEN="\e[32m"
YELLOW="\e[33m"
RESET="\e[0m"
echo -ne "${YELLOW}=>${RESET} BUILD - "
if [ "$build_result" == "0" ]; then
    echo -e "${GREEN}SUCCEEDED${RESET}"
else
    echo -e "${RED}FAILED${RESET}"
fi

make test
