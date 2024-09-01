package format

type DataFormat interface {
	Marshal(data any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}
