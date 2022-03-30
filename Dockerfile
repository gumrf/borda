FROM golang:1.17-alpine as Builder

RUN  apk add --update make

COPY . /borda/
WORKDIR /borda

RUN go mod download

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o ./build/borda-backend ./cmd/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /borda

COPY --from=Builder /borda/build/* /borda/
COPY ./migrations/* /borda/migrations/

EXPOSE 8080

ENTRYPOINT ["./borda-backend"]