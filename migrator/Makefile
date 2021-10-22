DOCKER_IMAGE = migrator

.PHONY: postgres-up
postgres-up:
	docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=zyxwv -d postgres

.PHONY: postgres-down
postgres-down:
	docker stop postgres

.PHONY: build
build:
	@echo "Docker Build..."
	docker build -t $(DOCKER_IMAGE) .

.PHONY: create-db
create-db:
	docker exec -it postgres createdb --username=postgres test_migration

.PHONY: drop-db
drop-db:
	docker exec -it postgres dropdb --username=postgres test_migration

.PHONY: migrate-up
migrate-up:build
	docker run --network host $(DOCKER_IMAGE) \
	-path=/migrations/ \
	-database "postgres://postgres:zyxwv@localhost:5432/test_migration?sslmode=disable" up 2

.PHONY: migrate-down
migrate-down:
	docker run --network host $(DOCKER_IMAGE) \
	-path=/migrations/ \
	-database "postgres://postgres:zyxwv@localhost:5432/test_migration?sslmode=disable" down 2
