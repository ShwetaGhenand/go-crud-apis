DOCKER_IMAGE = user-apis

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
