run-stack:
	@docker compose up -d

stop-stack:
	@docker compose down

run:
	go run .

build:
	go build .