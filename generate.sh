#!/bin/bash

if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <package_name> <file_name>"
    exit 1
fi

PACKAGE_NAME=$1
FILE_NAME=$2

mkdir -p generated/${PACKAGE_NAME}api

oapi-codegen -generate types -package "${PACKAGE_NAME}api" "$FILE_NAME" > "generated/${PACKAGE_NAME}api/types.gen.go"
oapi-codegen -generate client -package "${PACKAGE_NAME}api" "$FILE_NAME" > "generated/${PACKAGE_NAME}api/client.gen.go"
