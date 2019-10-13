#!/bin/bash

set -e

# Setup
cd "$(dirname $0)"
ROOT=$(pwd)
rm -rf debroot|| echo "no debroot to delete. good."
TARGET=$ROOT/debroot/nokeval_reader/
READER_TARGET=$TARGET/opt/nokeval_reader/
SERVICE_TARGET=$TARGET/etc/systemd/system/
mkdir -p "$TARGET"
mkdir -p "$READER_TARGET"
mkdir -p "$SERVICE_TARGET"

# Start building
GOOS=linux GOARCH=arm GOARM=5 go build -o nokeval_reader_arm cmd/reader/main.go

# Package
cp nokeval_reader_arm "$READER_TARGET"/nokeval_reader
cp -r DEBIAN "$TARGET"
cp -r LICE* "$READER_TARGET"
cp -r READ* "$READER_TARGET"
cp -r config.yaml.example "$READER_TARGET"
cp nokeval_reader.service "$SERVICE_TARGET"

cd debroot

dpkg -b nokeval_reader

VERSION=$(dpkg -I nokeval_reader.deb |grep Version|sed -e 's/Version: //g' -e 's/ //g')

mv nokeval_reader.deb nokeval_reader_$VERSION.deb
