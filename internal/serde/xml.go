package serde

import (
	"encoding/xml"
	"io"
)

type XMLSerde struct{}

func NewXMLSerde() *XMLSerde {
	return &XMLSerde{}
}

func (h *XMLSerde) Serialize(w io.Writer, data any) error {
	return xml.NewEncoder(w).Encode(data)
}

func (h *XMLSerde) Deserialize(r io.Reader, v any) error {
	return xml.NewDecoder(r).Decode(v)
}
