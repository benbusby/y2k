#!/bin/sh

mkdir -p out/
rm -rf out/*

platforms=(
    "windows/386"
    "windows/amd64"
    "windows/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "linux/arm"
    "linux/amd64"
    "linux/arm64"
    "linux/386")

for platform in "${platforms[@]}"
do
    echo "Compiling for $platform..."
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}

    output_name="y2k-$GOOS-$GOARCH"
    if [ $GOOS = "darwin" ]; then
        output_name="y2k-macos-$GOARCH"
    elif [ $GOARCH = "arm" ]; then
        output_name="y2k-$GOOS-arm32"
    elif [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o out/$output_name $package
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done
