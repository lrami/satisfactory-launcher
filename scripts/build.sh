#!/bin/bash

FOLDER="dist"
BINARY_NAME="launcher.exe"

source scripts/clean.sh

go build -o $FOLDER/$BINARY_NAME

cp template/* dist

echo "--> BUILD DONE"