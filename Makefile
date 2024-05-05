run: build
	@./bin/redis-clone

build:
	@go build -o bin/redis-clone .