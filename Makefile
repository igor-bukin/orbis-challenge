export GO111MODULE=on

BIN_NAME := $(or $(PROJECT_NAME),'api')
PKG_PATH := $(or $(PKG),'.')
PKG_LIST := $(shell go list ${PKG_PATH}/src/... | grep -v /vendor/)

MIGRATE=migrate -path sql/migrations -database postgres://postgres:12345@localhost:5432/orbis-api-db?sslmode=disable
TEST_MIGRATE=migrate -path sql/migrations  -database postgres://postgres:12345@localhost:5432/orbis-api-db-test?sslmode=disable

VERSION='1.0.0'
COMMIT=`git rev-parse HEAD`
BUILD_TIME=`date +%FT%T%z`

LD_FLAGS="-X ${PKG}/src/services/stats.version=${VERSION} \
 -X ${PKG}/src/services/stats.buildTime=${BUILD_TIME} \
 -X ${PKG}/src/services/stats.commit=${COMMIT}"

.PHONY: all build clean citest test cover lint

all: build lint test

dep: # Download required dependencies
	go mod vendor

lint: dep ## Lint the files local env
	cd src/ && golangci-lint run -c ../.golangci.yml && cd ../

citest: dep ## Run unit tests ci env
	mkdir -p ${CI_PROJECT_DIR}/reports
	go test -json > reports/test-report.out -coverprofile=reports/coverage-report.out -short -count=1 ${PKG_LIST}
	cat reports/test-report.out

cilint: dep
	mkdir -p ${CI_PROJECT_DIR}/reports
	golangci-lint run --timeout=5m -c .golangci.yml ./src/... > reports/gometalinter-report.out || cat reports/gometalinter-report.out
	cat reports/gometalinter-report.out

test: dep test-db-prepare ## Run unit tests
	cd src/ && go test -race -count=1 -short ./... && cd ../

msan: ## Run memory sanitizer
	go test -msan -short ${PKG_LIST}

cover:
	go test $(shell go list ./... | grep -v /vendor/;) -cover -v

build: dep ## Build the binary file
	@cd src; \
	go build -o ../bin/${BIN_NAME} -a -tags netgo -ldflags '-w -extldflags "-static"' -ldflags ${LD_FLAGS} .

clean: ## Remove previous build
	rm -f bin/$(BIN_NAME)

run: ## run application
	go run src/main.go

migrate-macos-install: ## Install migration tool on MacOS
	brew install golang-migrate

migrate-linux-install: ## Install migration tool on Linux Debian
	curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | sudo bash
	sudo apt-get install migrate=4.4.0

migrate-create: ## Create migration file with name
	migrate create -ext sql -dir sql/migrations -seq -digits 3 goscore

migrate-up: ## Run migrations
	$(MIGRATE) up

migrate-down: ## Rollback migrations
	$(MIGRATE) down

migrate-fix: ## Fix migrations
	$(MIGRATE) force $(v)

test-db-prepare: ## cleanup test db
	docker exec -u postgres postgres dropdb orbis-api-db-test || true
	docker exec -u postgres postgres createdb orbis-api-db-test
	$(TEST_MIGRATE) up

dc-up: ## up dockerized infrastructure
	@docker-compose -f ./infrastructure/docker-compose.yml up -d

dc-stop: ## stop dockerized infrastructure
	@docker-compose -f ./infrastructure/docker-compose.yml stop

dc-clean: ## clean up dockerized infrastructure
	@cd ./infrastructure ; docker-compose stop ; docker-compose rm -f

dc-show: ## show docker containers
	@docker container ls --format "{{.Names}} [{{.Ports}}]"

dc-postgres:
	@docker exec -it postgres psql api-db postgres_temp

fmt: ## format source files
	go fmt github.com/orbis-challenge/src/...

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

