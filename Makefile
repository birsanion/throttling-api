.PHONY: all build clean

all: build

server: build

build:
	go build -o server

run: server
	./server

clean:
	@rm -f server

image:=throttling-api
docker_build:
	DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 \
	docker build -t "${image}" .

docker_clean:
	@docker rmi -f "${image}"