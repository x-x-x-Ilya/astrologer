upgrade:
	go get -u ./... && go mod tidy

lint:
	cd internal && golangci-lint run ./...

build-app:
	cd ./deployments && \
	docker compose -f docker-compose.build.yml build

up:
	cd ./deployments && \
	docker compose -f docker-compose.up.yml --env-file .env up -d
