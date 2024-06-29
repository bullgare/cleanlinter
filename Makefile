.PHONY: test
test:
	go test ./...

.PHONY: integration-test
integration-test:
	cd ./test && go test ./...