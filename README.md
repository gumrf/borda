# Borda backend

## Requirements

- Docker
- docker-compose
- Go
- Make

## Build images
    docker-compose build

## Run
    docker-compose up borda

## Debug
    docker-compose up db
    make run

## Run pgweb(Postgres GUI)
    export DATABASE_URL=postgres://postgres:postgres@db:5432/postgres?sslmode=disable
    docker-compose up pgweb