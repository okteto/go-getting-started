.PHONY: start
start:
	air main.go

.PHONY: debug
debug:
	dlv debug --headless --listen=:2345 --log --api-version=2

