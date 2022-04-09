@PHONY: bins
bins:
	go install github.com/vektra/mockery/v2@latest


@PHONY: gen-mocks
gen-mocks:
	$(shell go env GOPATH)/bin/mockery --name CloudTrailClient --dir ./internal/aws

test: out_dir
	go test -coverprofile=_out/coverage.out ./...
	go tool cover -html _out/coverage.out

out_dir:
	@mkdir _out || true

build: out_dir
	GOOS=linux go build -o _out/main main.go