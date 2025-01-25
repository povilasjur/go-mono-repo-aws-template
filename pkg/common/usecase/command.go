package usecase

import (
	"context"
	"log"
	"time"
)

type CommandHandler[C any, R any] interface {
	Handle(ctx context.Context, cmd C) (R, error)
}

// CommandDecorator is a wrapper around Command that adds additional behavior
type CommandDecorator[C any, R any] struct {
	base CommandHandler[C, R]
}

// Execute logs before and after executing the wrapped command
func (d CommandDecorator[C, R]) Handle(ctx context.Context, cmd C) (R, error) {
	start := time.Now()
	log.Printf("Starting execution of command: %s", cmd)

	result, err := d.base.Handle(ctx, cmd)

	log.Printf("Finished execution in %s with response: %s", time.Since(start), result)

	return result, err
}
