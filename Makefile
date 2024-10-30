.PHONY: migrate migrate_down migrate_up migrate_version


# ==============================================================================
# Go migrate postgresql

force:
	migrate -database postgres://postgres:postgres@localhost:5432/weather4you?sslmode=disable -path migrations force 1

version:
	migrate -database postgres://postgres:postgres@localhost:5432/weather4you?sslmode=disable -path migrations version

migrate_up:
	migrate -database postgres://postgres:postgres@localhost:5432/weather4you?sslmode=disable -path migrations up 1

migrate_down:
	migrate -database postgres://postgres:postgres@localhost:5432/weather4you?sslmode=disable -path migrations down 1


# ==============================================================================
# Go run project

run:
	go run cmd/api/main.go

fillup:
	go run cmd/fillup/main.go	

build:
	go build ./cmd/api/main.go

test:
	go test -cover ./...


# ==============================================================================
# Docker compose

docker_delve:
	echo "Starting docker debug environment"
	sudo docker compose -f docker-compose.delve.yml up --build

docker_local_up:
	sudo docker compose -f docker-compose.yml up -d

local_build:
	echo "Starting local environment"
	sudo docker compose -f docker-compose.yml up --build

docker_down:
	sudo docker compose down

docker_build:
	sudo docker compose build

# ==============================================================================
# Tools commands

run-linter:
	echo "Starting linters"
	golangci-lint run ./...

swaggo:
	echo "Starting swagger generating"
	swag init -g **/**/*.go