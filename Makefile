.DEFAULT_GOAL := run
GO_BUILDER_IMAGE := golang:1.17.7-alpine
IMAGE := alpine:3.15.0

dep:
	go mod download

build-local:
	go build -o ./build/borda-backend

build-docker:
	docker build
	
run:
	go run cmd/main.go	