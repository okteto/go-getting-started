.PHONY: start
start:
	go run main.go

.PHONY: debug
debug:
	dlv debug --headless --listen=:2345 --log --api-version=2

.PHONY: push
push:
	okteto build -t okteto/hello-world:golang-dev --target dev .
	okteto build -t okteto/hello-world:golang .
