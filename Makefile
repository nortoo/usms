.PHONY: protobuf server

protobuf:
	@echo "Compiling protobuf files..."
	@bash hack/protobuf/compile.sh pkg/proto
	@git add .

server:
	@echo "Building docker image for server..."
	@bash build/server/build.sh
