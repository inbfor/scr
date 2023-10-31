lint:
	@echo "Running linter"
	go mod vendor
	docker run --rm -v $(PWD):/work:ro -w /work -it golangci/golangci-lint:latest golangci-lint run -v;