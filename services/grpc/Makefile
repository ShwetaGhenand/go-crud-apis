DOCKER_IMAGE = grpc-apis

MIGRATE_SCHEMA_NAME = create_users_table

.PHONY: create-migrate-schema
create-migrate-schema:
	migrate create -ext sql -dir users/sql/migration -seq $(MIGRATE_SCHEMA_NAME)


.PHONY: protoc
protoc:
	protoc --proto_path=proto proto/*.proto --go_out=proto/gen
	protoc --proto_path=proto proto/*.proto --go-grpc_out=proto/gen

.PHONY: build
build:
	@echo "Docker Build..."
	docker build -t $(DOCKER_IMAGE) .

.PHONY: up
up:
	docker-compose up -d

.PHONY: logs
logs:
	docker-compose logs --follow

.PHONY: down
down:
	docker-compose down

.PHONY: test
test:
	go test users/endpoint_test.go

.PHONY: buf-generate
buf-generate:
	buf generate
