.DEFAULT_GOAL := run
dep:
	go mod download

build: dep
	go build -o ./.bin/borda

run:
	go run cmd/main.go
	