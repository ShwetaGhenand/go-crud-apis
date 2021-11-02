DOCKER_IMAGE = user-apis
MIGRATOR_DOCKER_IMAGE = migrator

.PHONY: build
build:
	@echo "Docker Build..."
	docker build -t $(DOCKER_IMAGE) .

.PHONY: up
up: build
	docker-compose up -d

.PHONY: logs
logs:
	docker-compose logs --follow

.PHONY: down
down:
	docker-compose down

.PHONY: build-migrator
build-migrator:
	@echo "Docker Build..."
	docker build -t $(MIGRATOR_DOCKER_IMAGE) cmds/server/migrator

.PHONY: migrate-up
migrate-up:build-migrator
	docker run --network host $(MIGRATOR_DOCKER_IMAGE) \
	-path=/migrations/ \
	-database "postgres://postgres:zyxwv@localhost:5432/userdb?sslmode=disable" up 2

.PHONY: migrate-down
migrate-down:build-migrator
	docker run --network host $(MIGRATOR_DOCKER_IMAGE) \
	-path=/migrations/ \
	-database "postgres://postgres:zyxwv@localhost:5432/userdb?sslmode=disable" down 2
