#!/bin/bash

FOLDER="dist"
ARCHIVE_NAME="satisfactory-launcher.zip"

source scripts/build.sh

(cd $FOLDER; 7z a $ARCHIVE_NAME *)

echo "--> PACKAGE DONE"