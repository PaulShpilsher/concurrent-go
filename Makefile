.PHONY: test
test: test-sync-runner test-channel-runner

.PHONY: test-channel-runner
test-channel-runner:
	go test -v ./concurrency/chan/runner/runner_test.go

.PHONY: test-sync-runner
test-sync-runner:
	go test -v ./concurrency/sync/runner/runner_test.go


.PHONY: run-example
run-example:
	go run ./examples/main.go -kind=sync -limit=10 -tasks=500
