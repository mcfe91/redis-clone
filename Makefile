run: build
	@./bin/redis-clone --listenAddr :5001

build:
	@go build -o bin/redis-clone .