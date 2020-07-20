#!/bin/sh
mkdir -p coverage

echo "running golang tests"

go test -v -coverprofile ./coverage/cover.out
go tool cover -html=./coverage/cover.out -o ./coverage/cover.html
