.PHONY: build
build:
	go build -v ./cmd/authservice

.DEFAULT_GOAL := build