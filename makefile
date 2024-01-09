BINARY_NAME=companyXYZ

run:
	@echo "Running server"
	@cd bin/; ./${BINARY_NAME}

tidy:
	@echo "downloading application dependencies"
	@go mod tidy

build:
	@echo "Building binary"
	@go build -o bin/${BINARY_NAME} cmd/companyXYZ/main.go

test:
	@echo "Running all test"
	@go test ./... -v

swagger:
	@echo "Generating swagger"
	@swag init -d cmd/companyXYZ/,http/
	@swag fmt --exclude internal/,templ/,env/


instuction:
	@echo "Now add your .env to the bin/ created bin folder and run 'make dev'"

setup:tidy swagger build instuction

dev: setup run

prod:
	@echo "YOUR PRODUCTION ENVIROMMENT MATTERS FIRST"

