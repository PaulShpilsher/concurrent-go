.PHONY: test
test:
	go test -v ./...

.PHONY: bench
bench:
	go test -bench . ./...

.PHONY: run-example-channel
run-example-channel:
	go run ./examples/chan/main.go


.PHONY: run-example-sync
run-example-sync:
	go run ./examples/sync/main.go
