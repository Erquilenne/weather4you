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


# ==============================================================================
# Docker compose

docker_up:
	sudo docker compose up -d


docker_down:
	sudo docker compose down

docker_build:
	sudo docker compose build