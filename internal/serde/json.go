package serde

import (
	"encoding/json"
	"io"
)

type JSONSerde struct{}

func NewJSONSerde() *JSONSerde {
	return &JSONSerde{}
}

func (h *JSONSerde) Serialize(w io.Writer, data any) error {
	return json.NewEncoder(w).Encode(data)
}

func (h *JSONSerde) Deserialize(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}
