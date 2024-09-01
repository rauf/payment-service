package protocol

import (
	"context"
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
		return nil, err
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			slog.ErrorContext(ctx, "Failed to close connection", err)
			return
		}
	}(conn)

	_, err = conn.Write(data)
	if err != nil {
		return nil, err
	}

	response := make([]byte, 4096)
	n, err := conn.Read(response)
	if err != nil {
		return nil, err
	}

	return response[:n], nil
}
