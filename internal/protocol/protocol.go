package protocol

import "context"

type Handler interface {
	Send(ctx context.Context, data []byte) ([]byte, error)
}
