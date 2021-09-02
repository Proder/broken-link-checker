run:
	docker-compose up  --remove-orphans --build

run_checker:
	go run -race cmd/checker/main.go

run_server_test:
	go run cmd/server_test/main.go

lint:
	gofumpt -w -s ./broken-link-checker/..
	golangci-lint run --fix

