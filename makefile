BINARY_NAME=BlockRide

run:
	@cd bin/; ./${BINARY_NAME}

compile:
	@echo "building binary"
	@go build -o bin/${BINARY_NAME} cmd/blockride/main.go

test:
	@echo "Running all test"
	@go test ./... -v

start: build start
