GO_BUILDER_IMAGE := golang:1.17.7-alpine
IMAGE := alpine:3.15.0

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/borda-backend ./cmd/main.go

run:
	go run ./cmd/borda-api-server/main.go serve

swag:
	swag init -g internal/app/app.go
	swag fmt -g internal/app/app.go

clean:
	rm -rf ./build