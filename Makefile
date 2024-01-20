all: build up

build:
	env GOOS=linux GOARCH=amd64 go build -o docker/frontend/frontend src/services/frontend/frontend.go
	env GOOS=linux GOARCH=amd64 go build -o docker/worker/worker src/services/worker/worker.go

up: build
	docker-compose up --build

down:
	docker-compose down