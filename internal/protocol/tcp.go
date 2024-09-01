package protocol

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
)

type TCPProtocol struct {
	Address string
}

func NewTCPConnection(address string) *TCPProtocol {
	return &TCPProtocol{
		Address: address,
	}
}

func (t *TCPProtocol) Send(ctx context.Context, data []byte) ([]byte, error) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", t.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to establish TCP connection: %w", err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			slog.ErrorContext(ctx, "Failed to close connection", "error", err)
		}
	}(conn)

	_, err = conn.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to write data to TCP connection: %w", err)
	}

	response := make([]byte, 4096)
	n, err := conn.Read(response)
	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("connection closed by remote host: %w", err)
		}
		return nil, fmt.Errorf("failed to read response from TCP connection: %w", err)
	}

	return response[:n], nil
}
