.PHONY: fmt
fmt:
	gofmt -w -s .
	goimports -w .
	go clean ./...
