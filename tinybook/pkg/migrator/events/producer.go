package events

import "context"

type Producer interface {
	ProduceInconsistentEvent(ctx context.Context, event InconsistentEvent) error
}
