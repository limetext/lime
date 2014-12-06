#!/usr/bin/env bash

source "$(dirname -- "$0")/warn.sh"
source "$(dirname -- "$0")/setup.sh"

ret=0

fold_start "fmt" "check formatting"
diff_test "go run tasks/build/fix.go"
let ret=$ret+$test_result
fold_end "fmt"

exit $ret
