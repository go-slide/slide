#!/bin/sh
if [ "$(go fmt | wc -l)" -gt 0 ]; then
    echo "exit"
    exit 1
fi