package concurrency

import "errors"

var (
	ErrNilArgument   = errors.New("nil argument")
	ErrRunnerClosed  = errors.New("runner closed")
	ErrChannelClosed = errors.New("channel closed")
)
