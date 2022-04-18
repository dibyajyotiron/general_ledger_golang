#!make

# Integration test location provided here
INTEGRATION_TEST_PATH?=./tests/integration-test

# set of env variables that you need for testing
ENV_LOCAL_TEST=\
  APP_ENV=prod\
  RUN_MODE=release\
  DB_TYPE=postgres\
  DB_USER=postgres\
  DB_PASSWORD=admin\
  DB_HOST=127.0.0.1\
  DB_PORT=5432\
  DB_NAME=golang_ledger\
  DB_TABLE_PREFIX=golang_ledger_\
  DB_SSL_MODE=disable\
  APP_PORT=3000\
  JWT_SECRET=usx1957-213123123123-12312sa7687-23424\

# Include .env
include .env
#export $(shell sed 's/=.*//' env)

# this command will start docker components that we set in docker-compose.yml
docker.start.components:
	docker-compose up -d --remove-orphans postgres;

# shutting down docker components
docker.stop:
	docker-compose down;

# this command will trigger integration test
# INTEGRATION_TEST_SUITE_PATH is used to run specific tests in Golang, if it's not specified
# it will run all tests under ./tests/integration-test directory
test.integration:
	$(ENV_LOCAL_TEST) \
	go test -tags=integration $(INTEGRATION_TEST_PATH) -count=1 -run=$(INTEGRATION_TEST_SUITE_PATH)

# this command will trigger integration test with verbose mode
test.integration.debug:
	$(ENV_LOCAL_TEST) \
	go test -tags=integration $(INTEGRATION_TEST_PATH) -count=1 -v -run=$(INTEGRATION_TEST_SUITE_PATH)