#!/bin/bash

set -e

outputPath="pkg/assets/cinny"

assetJSON="$(curl -s https://api.github.com/repos/ajbura/cinny/releases/latest | jq '.assets[] | select(.name | test(".tar.gz"))')"

name="$(echo $assetJSON | jq -r '.name')"
url="$(echo $assetJSON | jq -r '.browser_download_url')"

echo Downloading $name...

curl -s -L "$url" | tar -C "$outputPath" --strip-components 1 -xzf -

echo Success!

echo Source files stored in $(realpath $outputPath)
