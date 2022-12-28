#!/bin/sh

SCRIPT_DIR="$(CDPATH= command cd -- "$(dirname -- "$0")" && pwd -P)"

echo "- Building executable"
go build

echo "- Running tests"
for dir in examples/*; do
    cd "$SCRIPT_DIR"

    if [ ! -d "$dir" ]; then
        continue
    fi

    cd "$dir"
    ./init.sh

    cd "$SCRIPT_DIR"

    # Test output with input arg 15
    # (needed for n-terms fibonacci example)
    output="$(./y2k "$dir" 15)"
    expected="$(cat "$dir/test-output.txt")"

    if [ "$output" != "$expected" ]; then
        echo "ERROR: $dir"
        echo "Expected: $expected"
        echo "Output: $output"
        exit 1
    else
        echo "OK: $dir"
    fi
done

echo "All tests passed"
