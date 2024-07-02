include .env

gen:
	@protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    protos/projectrol.proto

migrations/new:
	@echo "Creating migration files for ${name}..."
	migrate create -seq -ext=sql -dir=migrations ${name}

migrations/up:
	@echo "Start migration up..."
	migrate -path=migrations -database=postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME} up

migrations/down:
	@echo "Start migration up..."
	migrate -path=migrations -database=postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME} down