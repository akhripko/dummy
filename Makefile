SHELL=/bin/bash
ROOT_DIR := $(shell pwd)
IMAGE_TAG := $(shell git rev-parse --short HEAD)
IMAGE_NAME := company/srv
REGISTRY := change-it.dkr.ecr.us-west-2.amazonaws.com

.PHONY: grpcgen gqlgen mockgen build run lint test test_integration dockerise deploy run_postgresql run_redis start_deps stop_deps

ci: mod lint build test dockerise

mod:
	go mod download
	go mod vendor

grpcgen:
	protoc -I api api/service.proto --go_out=plugins=grpc:api

gqlgen:
	cd src/srv/srvgql && \
	rm -f generated.go models/*_gen.go && \
	go run scripts/gqlgen.go -v

build:
	go build -o artifacts/svc ./cmd/svc/main.go

run:
	go run ./cmd/svc/main.go

lint:
	cd ./src && golangci-lint run

mockgen:
	mockgen -source=src/service/service.go -destination=src/service/mock/deps.go
	mockgen -source=src/srv/srvhttp/service.go -destination=src/srv/srvhttp/mock/service.go
	mockgen -source=src/srv/srvgrpc/service.go -destination=src/srv/srvgrpc/mock/service.go
	mockgen -source=src/srv/srvgql/service.go -destination=src/srv/srvgql/mock/service.go

test:
	go test -cover -v `go list ./src/...`

test_integration:
	INTEGRATION_TEST=YES go test -cover -v `go list ./...`

dockerise:
	docker build -t ${IMAGE_NAME}:${IMAGE_TAG} -f ./cmd/svc/Dockerfile .
	docker tag ${IMAGE_NAME}:${IMAGE_TAG} ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}

deploy:
	`AWS_SHARED_CREDENTIALS_FILE=~/.aws/credentials AWS_PROFILE=xid aws ecr get-login --region us-west-2 --no-include-email`
	docker push ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}
	#docker tag ${IMAGE_NAME}:${IMAGE_TAG} ${REGISTRY}/${IMAGE_NAME}:latest
	#docker push ${REGISTRY}/${IMAGE_NAME}:latest

run_postgresql:
	docker run -d --name dummy_postgresql -e POSTGRES_DB=dummy -v ${ROOT_DIR}/tmp/sql/data:/var/lib/postgresql/data -p 5432:5432 postgres:11

run_redis:
	docker run --name dummy_redis -p 6379:6379 -d redis

start_deps:
	docker start dummy_redis
	docker start dummy_postgresql

stop_deps:
	docker stop dummy_redis
	docker stop dummy_postgresql

#.PHONY: exec_redis_sh
#exec_redis_sh:
#	docker exec -it dummy_redis sh
#    # redis-cli