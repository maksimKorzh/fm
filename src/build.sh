#!/bin/bash

export GOOS=linux
go build fm.go
export GOOS=windows
go build fm.go