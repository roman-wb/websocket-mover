.PHONY: server ngrok lint build docker-build docker-run

server:
	go run main.go

ngrok:
	go run main.go &
	ngrok http 8080
	killall -9 main

lint:
	docker run --rm -v `pwd`:/app -w /app golangci/golangci-lint golangci-lint run -v

build:
	go build -ldflags "-s -w" -o bin/server .

docker-build:
	docker build -t websocket-mover .

docker-run:
	docker run --rm -it -p 8080:8080 websocket-mover
