run:
	go run cmds/registrar.go --postgres-host localhost --postgres-port 5432 --postgres-db postgres --postgres-user postgres --postgres-password password --domain localhost --server-port 8080

postgres:
	docker run --name postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=postgres -p 5432:5432 -d postgres

stop-postgres:
	docker stop postgres && docker rm postgres

build: ## Bulil the server
	go build -o bin/server cmds/registrar.go

server-start:
	@go run cmds/registrar.go \
		--server-port 8080 \
		--debug \
		--domain localhost \
		--sql-log-level 4 \
		--postgres-host localhost \
		--postgres-port 5432 \
		--postgres-db postgres \
		--postgres-user postgres \
		--postgres-password password \

getverifiers:
	@echo "Installing staticcheck" && go get -u honnef.co/go/tools/cmd/staticcheck && go install honnef.co/go/tools/cmd/staticcheck
	@echo "Installing gocyclo" && go get -u github.com/fzipp/gocyclo/cmd/gocyclo && go install github.com/fzipp/gocyclo/cmd/gocyclo
	@echo "Installing deadcode" && go get -u github.com/remyoudompheng/go-misc/deadcode && go install github.com/remyoudompheng/go-misc/deadcode
	@echo "Installing misspell" && go get -u github.com/client9/misspell/cmd/misspell && go install github.com/client9/misspell/cmd/misspell
	@echo "Installing golangci-lint" && go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1
	go mod tidy

verifiers: fmt lint cyclo deadcode spelling staticcheck

checks: verifiers

fmt:
	@echo "Running $@"
	@gofmt -d .

lint:
	@echo "Running $@"
	@$(shell go env GOPATH)/bin/golangci-lint run

cyclo:
	@echo "Running $@"
	@$(shell go env GOPATH)/bin/gocyclo -over 100 .

deadcode:
	@echo "Running $@"
	@$(shell go env GOPATH)/bin/deadcode -test $(shell go list ./...) || true

spelling:
	@echo "Running $@"
	@$(shell go env GOPATH)/bin/misspell -i monitord -error `find .`

staticcheck:
	@echo "Running $@"
	@$(shell go env GOPATH)/bin/staticcheck -- ./...

