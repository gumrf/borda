# Borda backend

## Requirements

- Docker
- docker-compose
- Go
- Make

## Set env variables
    SERVER_ADDR="0.0.0.0:8080"

    POSTGRES_HOST=db
    POSTGRES_USER=postgres
    POSTGRES_PASSWORD=postgres
    POSTGRES_DB=postgres

    DATABASE_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:5432/${POSTGRES_DB}?sslmode=disable"
    GITHUB_ACCESS_TOKEN=<token>
    TASK_REPOSITORY_URL=<url>

## Build images
    docker-compose build

## Run
    docker-compose up borda

## Debug
    docker-compose up db
    make run

## Run pgweb(Postgres GUI)
    docker-compose up pgweb