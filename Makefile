.PHONY: build
build:
	go build -v ./cmd/authservice
	./authservice.exe

.DEFAULT_GOAL := build

docker:
	docker-compose -f docker-compose.yml up

run: 
	docker run -d -p 8080:8080 auth_service

stop:
	docker stop auth_service-container || echo "контейнер не запущен"