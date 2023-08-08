build:
	docker compose build --no-cache

up:
	docker-compose down
	DOCKER_BUILDKIT=0 COMPOSE_DOCKER_CLI_BUILD=0 docker compose up 

test:
	go test -v ./... -coverpkg ./...