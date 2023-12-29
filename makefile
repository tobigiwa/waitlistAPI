BINARY_NAME=BlockRide

run:
	@echo "Running server"
	@cd bin/; ./${BINARY_NAME}

tidy:
	@echo "downloading application dependencies"
	@go mod tidy

build:
	@echo "Building binary"
	@go build -o bin/${BINARY_NAME} cmd/blockride/main.go

test:
	@echo "Running all test"
	@go test ./... -v

swagger:
	@echo "Generating swagger"
	@swag init -d cmd/blockride/,http/
	@swag fmt --exclude internal/,templ/,env/


instuction:
	@echo "Now add your .env to the bin/ created bin folder and run 'make dev'"

setup:tidy swagger build instuction

dev: run

prod: tidy build

