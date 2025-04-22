install:
	@echo "Downloading dependecies..."
	@go get
	@echo "Validating dependecies..."
	@go mod tidy
	@echo "Creating vendor..."
	@go mod vendor
	@echo "Installation completed successfully."

build:
	@echo "Building project..."
	@go build
	@echo "Build completed successfully."

run-worker:
	@echo "Running application..."
	@go run main.go -worker

run-script::
	@echo "Running application..."
	@go run main.go -script

test:
	@echo "Running project tests..."
	@go test -v -cover ./...
	@echo "Running project tests..."

lint:
	@echo "Running golangci-lint..."
	@golangci-lint run
	@echo "Linter completed successfully. No issues found."

docker-setup:
	@echo "Starting docker services..."
	@docker-compose up -d

coverage:
	@echo "Running project coverage..."
	@go test ./... -coverprofile fmtcoverage.html fmt
	@go test ./... -coverprofile cover.out
	@go tool cover -html=cover.out
	@go tool cover -html=cover.out -o cover.html
	@echo "Coverage completed successfully."
