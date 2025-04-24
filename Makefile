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

run-api::
	@echo "Running application..."
	@go run main.go -api

test:
	@echo "Running project tests..."
	@go test -v -cover ./...
	@echo "Running project tests..."

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

generate-mocks:
	@echo "Gerando mocks..."
	@mockgen -source=application/usecases/pulse/pulse.go -destination=application/usecases/mocks/mock_pulse_usecase.go -package=Usecasemocks
	@mockgen -source=infrastructure/repository/pulse/pulse.go -destination=infrastructure/repository//mocks/mock_pulse_repository.go -package=repositorymocks
	@mockgen -source=cmd/api/task/populate_queue_task/populate_queue_task.go -destination=cmd/api/task/populate_queue_task/mocks/mock_populate_task.go -package=populateTask
	@echo "Mocks atualizados com sucesso!"