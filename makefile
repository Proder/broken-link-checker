run_checker:
	go run -race cmd/checker/main.go

lint:
	gofumpt -w -s ./broken-link-checker/..
	golangci-lint run --fix

