run:
	docker-compose up --build --remove-orphans 

postgres-up:
	migrate -source file:./migrations -database postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable up

postgres-down:
	migrate -source file:./migrations -database postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable down

lint:
	golangci-lint run