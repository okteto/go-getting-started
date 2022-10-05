.PHONY: start
start:
	go run main.go

.PHONY: test
test: 
	go test ./pkg/...
 
.PHONY: debug
debug:
	dlv debug --headless --listen=:2345 --log --api-version=2

