package format

import (
	"encoding/json"
)

type JSON struct {
	URL string
}

func NewJSONProtocol() *JSON {
	return &JSON{}
}

func (h *JSON) Marshal(data any) ([]byte, error) {
	return json.Marshal(data)
}

func (h *JSON) Unmarshal(data []byte, v any) error {
	err := json.Unmarshal(data, &v)
	return err
}
