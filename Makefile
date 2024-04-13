ifeq ($(wildcard .env), .env)
include .env
export
endif

check_env_file:
	@if [ -f .env ]; then \
		echo ".env file exists"; \
	else \
		echo "Error: .env file not found"; \
		exit 1; \
	fi

TEST_APP_PORT=8888
TEST_PG_PORT=4444
TEST_RMQ_PORT=5556
TEST_RMQ_MANAGEMENT_PORT=55556

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
	-docker compose -f ./deployments/service-compose.yml -p banner down

test-env-file:
	echo POSTGRES_HOST=test_db > ./.test.env
	echo POSTGRES_PORT=$(TEST_PG_PORT) >> ./.test.env
	echo POSTGRES_USER=test_user >> ./.test.env
	echo POSTGRES_PASSWORD=crakme >> ./.test.env
	echo POSTGRES_DB=banner >> ./.test.env

	echo RABBITMQ_HOST=test_rmq >> ./.test.env
	echo RABBITMQ_PORT=$(TEST_RMQ_PORT) >> ./.test.env
	echo RABBITMQ_PORT_MANAGEMENT=$(TEST_RMQ_MANAGEMENT_PORT) >> ./.test.env

	echo APP_PORT=$(TEST_APP_PORT) >> ./.test.env
	echo ADMIN_TOKEN=admin >> ./.test.env
	echo USER_TOKEN=user >> ./.test.env
	echo PG_URL=postgres://test_user:crakme@test_db:5432/banner >> ./.test.env
	echo RMQ_URL=amqp://guest:guest@test_rmq:5672/ >> ./.test.env

test-env: build-image test-env-file
	$(eval include .test.env)
	$(eval export)
	$(eval export PG_URL=postgres://test_user:crakme@localhost:4444/banner)

	docker compose -f ./deployments/test-env.yml -p test-banner up -d

test: test-env
	-go test -v ./...
	make clean-test-env

worker:
	go run cmd/worker/main.go

clean-test-env:
	-$(eval include .test.env)
	-$(eval export)
	-$(eval export PG_URL=postgres://test_user:crakme@localhost:4444/banner)
	-docker compose -f ./deployments/test-env.yml -p test-banner down
	-docker volume rm test-banner_test-pg-data
	-docker volume rm test-banner_test-rabbitmq-data
	-rm ./.test.env

clean-prod-env: stop
	-rm ./.env
	-docker volume rm banner_postgres_data
	-docker volume rm banner_rabbitmq-data

clean-docker:
	-docker system prune -a

clean: stop clean-test-env clean-prod-env clean-docker
	echo "Successful cleaned ;3"