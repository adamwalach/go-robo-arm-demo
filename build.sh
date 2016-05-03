#!/bin/bash

GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -o=main *.go
