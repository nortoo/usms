.PHONY: protobuf server test

protobuf:
	@echo "Compiling protobuf files..."
	@bash hack/protobuf/compile.sh pkg/proto
	@git add .

server:
	@echo "Building docker image for server..."
	@bash build/server/build.sh

test:
	@echo "Running tests..."
	@go test -v ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@rm coverage.out