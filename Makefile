run: build 
	@./bin/dpoker

build: 
	@go build -o bin/dpoker

test: 
	@go test -v ./...
