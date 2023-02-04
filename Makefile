.PHONY: test
test:
	go test -v ./concurrent/runner_test.go 

.PHONY: test
run:
	go run main.go
