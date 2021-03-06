#!/bin/sh
if [ "$(go fmt | wc -l)" -gt 0 ]; then
    echo "run go fmt before PR"
    exit 1
fi

if [ "$(golint | wc -l)" -gt 0 ]; then
    echo "run go lint before PR"
    exit 1
fi