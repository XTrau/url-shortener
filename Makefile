run:
	CGO_ENABLED=0 go run ./cmd/app
up:
	sudo docker compose down
	sudo docker compose build
	sudo docker compose up -d
format:
	go fmt ./...
	go vet ./...
