FROM golang:1.17-alpine as Builder

RUN  apk add --update make

ADD . /borda-backend/
WORKDIR /borda-backend

RUN go mod download
RUN CGO_ENABLED=0 go build -o ./build/borda-backend

# 
FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=Builder /borda-backend/build/* /opt/Borda/

EXPOSE 6969

ENTRYPOINT ["/opt/Borda/borda-backend"]