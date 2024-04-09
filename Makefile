include .env

export
ifndef $(SSL_MODE)
	export SSL_MODE=disable
endif

export PG_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(SSL_MODE)



build:
	go build cmd/banner/main.go

run:
	go run cmd/banner/main.go

migrate:
	go run ./cmd/migrator --config=./configs/local.yaml

test:
	env
	go test -v ./tests/...

gen:
	go generate -x ./...
env:
	cat ./.env.example > ./.env

db:
	docker-compose up -d --remove-orphans
