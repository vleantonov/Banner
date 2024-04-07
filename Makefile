build:
	go build cmd/banner/main.go

run:
	go run cmd/banner/main.go --config=./configs/local.yaml

migrate:
	go run ./cmd/migrator --config=./configs/local.yaml

test:
	go test -v ./tests/...

gen:
	go generate -x ./...
env:
	cat ./.env.example > ./.env

db:
	docker-compose up -d --remove-orphans
