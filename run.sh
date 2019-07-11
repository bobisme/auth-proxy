#!/bin/sh

echo "-s -g '*.go' -- go run main.go
-s -g 'example/*.go' -- go run example/main.go" \
  | reflex --decoration=fancy -c -
