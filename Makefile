# Docker commands
up:
	docker compose -f docker/docker-compose.yml up -d

down:
	docker compose -f docker/docker-compose.yml down

logs:
	docker compose -f docker/docker-compose.yml logs -f

ps:
	docker compose -f docker/docker-compose.yml ps

# App commands
build:
	go build -o bin/api cmd/api/main.go
	go build -o bin/consumer cmd/consumer/main.go

run-api:
	go run cmd/api/main.go
run-consumer:
	go run cmd/consumer/main.go

# Full start (api)
start: up build run-api


clean:
	docker compose -f docker/docker-compose.yml down -v
	rm -rf bin/