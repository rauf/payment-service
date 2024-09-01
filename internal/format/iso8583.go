package format

type ISO8583 struct {
}

func NewISO8583Protocol() *ISO8583 {
	return &ISO8583{}
}

func (h *ISO8583) Marshal(data any) ([]byte, error) {
	return nil, nil
}

func (h *ISO8583) Unmarshal(data []byte, v any) error {
	return nil
}
