run:
	docker-compose up --build --remove-orphans 

postgres-up:
	migrate -source file:./migrations/postgres -database postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable up

postgres-down:
	migrate -source file:./migrations/postgres -database postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable down


clickhouse-up:
	goose -dir ./migrations/clickhouse clickhouse "tcp://127.0.0.1:19000" up

clickhouse-down:
	goose -dir ./migrations/clickhouse clickhouse "tcp://127.0.0.1:19000" down

lint:
	golangci-lint run