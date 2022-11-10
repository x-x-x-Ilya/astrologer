upgrade-backend:
	go get -u ./... && go mod tidy

lint-backend:
	cd internal && golangci-lint run ./...