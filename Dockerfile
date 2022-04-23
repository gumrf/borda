FROM golang:1.17-alpine as builder

RUN  apk add --update make

COPY . /borda/
WORKDIR /borda

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o ./bin/borda-api-server ./cmd/borda-api-server/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /borda

COPY --from=builder /borda/bin/* /borda/
COPY ./migrations/* /borda/migrations/

EXPOSE 8080

ENTRYPOINT ["./borda-api-server"]
CMD ["serve"]
# CMD ["serve", "--env", "production", "--hostname", "0.0.0.0", "--port", "8080"]