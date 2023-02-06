.PHONY: test
test:
	go test -v ./runner/runner_test.go 

.PHONY: test
run-example:
	go run ./examples/main.go -limit=10 -tasks=500
