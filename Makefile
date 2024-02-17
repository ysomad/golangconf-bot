include .env
export

export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING=${PG_URL}

.PHONY: lint
lint:
	golangci-lint run

.PHONY:
run:
	go mod tidy && go mod download && go run .

.PHONY: test
test:
	go test -v -cover -race -count 1 -coverpkg=$(go list ./... | grep -v '/mock$' | tr '\n' ',') ./internal/...

.PHONY: mock
mock:
	rm -rf ./internal/bot/mock/ \
	mockery

.PHONY: compose-up
compose-up:
	docker-compose up --build -d && docker-compose logs -f

.PHONY: compose-down
compose-down:
	docker-compose down --remove-orphans

.PHONY: goose-new
goose-new:
	@read -p "Enter the name of the new migration: " name; \
	goose -dir migrations create $${name// /_} sql

.PHONY: goose-up
goose-up:
	@echo "Running all new database migrations..."
	goose -dir migrations validate
	goose -dir migrations up

.PHONY: goose-down
goose-down:
	@echo "Running all down database migrations..."
	goose -dir migrations down

.PHONY: goose-reset
goose-reset:
	@echo "Dropping everything in database..."
	goose -dir migrations reset

.PHONY: goose-status
goose-status:
	goose -dir migrations status
