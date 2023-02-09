.PHONY: test
test: test-sync-runner test-channel-runner

.PHONY: test-channel-runner
test-channel-runner:
	go test -v ./concurrency/chan/runner/runner_test.go

.PHONY: test-sync-runner
test-sync-runner:
	go test -v ./concurrency/sync/runner/runner_test.go


.PHONY: run-example-channel
run-example-channel:
	go run ./examples/chan/main.go


.PHONY: run-example-sync
run-example-sync:
	go run ./examples/sync/main.go
