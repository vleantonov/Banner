ifeq ($(wildcard .env), .env)
include .env
export
export PG_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)
endif

check_env_file:
	@if [ -f .env ]; then \
		echo ".env file exists"; \
		env; \
	else \
		echo "Error: .env file not found"; \
		exit 1; \
	fi

TEST_APP_PORT=8888
TEST_PG_PORT=4444

.DEFAULT_GOAL := run

gen:
	go generate -x ./...

build-image:
	docker image build -f ./build/Dockerfile -t banner_app .

env:
	cat ./.env.example > ./.env

run: check_env_file build-image

	docker compose -f ./deployments/service-compose.yml -p banner up -d --remove-orphans
stop:
	docker compose -f ./deployments/service-compose.yml -p banner down

run-local: check_env_file
	go run cmd/migrator/main.go
	go run cmd/banner/main.go

test-env-file:
	echo POSTGRES_PORT=$(TEST_PG_PORT) > ./.test.env
	echo POSTGRES_USER=test_user >> ./.test.env
	echo POSTGRES_PASSWORD=crakme >> ./.test.env
	echo POSTGRES_DB=banner >> ./.test.env
	echo APP_PORT=$(TEST_APP_PORT) >> ./.test.env
	echo ADMIN_TOKEN=admin >> ./.test.env
	echo USER_TOKEN=user >> ./.test.env
	echo PG_URL=postgres://test_user:crakme@test_db:5432/banner >> ./.test.env
test-env: build-image test-env-file
	$(eval include .test.env)
	$(eval export)
	$(eval export PG_URL=postgres://test_user:crakme@localhost:4444/banner)

	docker compose -f ./deployments/test-env.yml -p test-banner up -d
test: test-env
	-go test -v ./...
	make clean-test-env
clean-test-env:
	$(eval include .test.env)
	$(eval export)
	$(eval export PG_URL=postgres://test_user:crakme@localhost:4444/banner)
	docker compose -f ./deployments/test-env.yml -p test-banner down
	docker volume rm test-banner_test-pg-data
	rm ./.test.env

clean-prod-env: stop
	-rm ./.env
	-docker volume rm banner_postgres_data
clean-docker:
	-docker system prune -a
clean: clean-test-env stop clean-prod-env clean-docker
	echo "Successful cleaned ;3"