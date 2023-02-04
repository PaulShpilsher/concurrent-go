.PHONY: test
test:
	go test -v ./runner/runner_test.go 

.PHONY: test
run:
	go run main.go
