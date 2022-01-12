fmt:
	find . -name '*.go' | grep -v vendor | xargs gofmt -s -w

depscheck:
	@echo "==> Checking source code with go mod tidy..."
	@go mod tidy
	@git diff --exit-code -- go.mod go.sum || \
		(echo; echo "Unexpected difference in go.mod/go.sum files. Run 'go mod tidy' command or revert any go.mod/go.sum changes and commit."; exit 1)
	@echo "==> Checking source code with go mod vendor..."
	@go mod vendor
	@git diff --compact-summary --ignore-space-at-eol --exit-code -- vendor || \
		(echo; echo "Unexpected difference in vendor/ directory. Run 'go mod vendor' command or revert any go.mod/go.sum/vendor changes and commit."; exit 1)

test:
	@go test ./...

lint:
	@echo "==> Checking source code against linters..."
	@if command -v golangci-lint; then (golangci-lint run ./...); else ($(GOPATH)/bin/golangci-lint run ./...); fi