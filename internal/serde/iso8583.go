package serde

import "io"

type ISO8583Serde struct {
}

func NewISO8583Serde() *ISO8583Serde {
	return &ISO8583Serde{}
}

func (h *ISO8583Serde) Serialize(w io.Writer, data any) error {
	return nil
}

func (h *ISO8583Serde) Deserialize(r io.Reader, v any) error {
	return nil
}
