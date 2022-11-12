upgrade-backend:
	go get -u ./... && go mod tidy

lint-backend:
	cd internal && golangci-lint run ./...

build:
	cd deployments/build && \
	docker-compose -f docker-compose.yml build

up:
	cd deployments/up && \
	docker-compose -f up/docker-compose.base.yml --env-file .env up -d
