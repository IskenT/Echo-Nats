run:
	docker-compose up  --remove-orphans --build

postgres-up:
	migrate -source file:./migrations -database postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable up

postgres-down:
	migrate -source file:./migrations -database postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable down