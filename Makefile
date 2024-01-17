all: build up

build:
	go build -o docker/frontend/frontend.exe src/services/frontend/frontend.go
	go build -o docker/worker/worker.exe src/services/worker/worker.go

up: build
	docker-compose up --build

down:
	docker-compose down