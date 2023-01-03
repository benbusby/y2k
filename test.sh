#!/bin/sh

SCRIPT_DIR="$(CDPATH= command cd -- "$(dirname -- "$0")" && pwd -P)"
TEST_DIR="$SCRIPT_DIR/test-output"

echo "- Building executable"
go build

echo "- Running tests"
for example in examples/*; do
    # Set up test directory for raw Y2K file exports
    rm -rf "$TEST_DIR"
    mkdir "$TEST_DIR"

    # Evaluate the expected output of a Y2K example file
    expected="$(./y2k $example 15)"

    # Export the raw file to a set of empty timestamp files
    ./y2k -outdir $TEST_DIR -export $example >/dev/null
    output="$(./y2k $TEST_DIR 15)"

    # Check if both outputs are equal
    if [ "$output" != "$expected" ]; then
        echo "ERROR: $example"
        echo "Expected: $expected"
        echo "Output: $output"
        exit 1
    else
        echo "OK: $example"
    fi
done

echo "All tests passed"
rm -rf "$TEST_DIR"
