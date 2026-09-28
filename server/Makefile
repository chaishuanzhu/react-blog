.PHONY: run build test lint sqlc up down import

ifneq (,$(wildcard .env))
include .env
export
endif

run:
	go run ./cmd/server

build:
	CGO_ENABLED=0 go build -o bin/blog-server ./cmd/server

test:
	go test ./...

lint:
	go vet ./...

sqlc:
	sqlc generate

# make import DIR=./cloudbase-export ARGS="-truncate -admin-email you@example.com"
import:
	go run ./cmd/import-cloudbase -dir $(DIR) $(ARGS)

# Builds and runs the single-container deployment defined in ../docker-compose.yml.
up:
	docker compose -f ../docker-compose.yml up -d --build

down:
	docker compose -f ../docker-compose.yml down
