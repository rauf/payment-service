package protocol

import "context"

// Handler is the interface for all protocol handlers.
type Handler interface {
	Send(ctx context.Context, data []byte) ([]byte, error)
}
