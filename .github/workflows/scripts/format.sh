#!/bin/sh
if [ "$(go fmt | wc -l)" -gt 0 ]; then
    echo "run go fmt before PR"
    exit 1
fi