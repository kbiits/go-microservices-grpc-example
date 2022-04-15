include .env

migrate-up:
	migrate -path ./migration/ -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:5432/${POSTGRES_DB}?sslmode=disable up

migrate-down:
	migrate -path ./migration/ -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:5432/${POSTGRES_DB}?sslmode=disable down

build:
	CGO_ENABLED=0 go build -o bin/main main.go

.PHONY: migrate test