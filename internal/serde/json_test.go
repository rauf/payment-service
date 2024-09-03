package serde

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

func TestJSONProtocol(t *testing.T) {
	jsonProtocol := NewJSONSerde()

	t.Run("Serialize", func(t *testing.T) {
		data := map[string]string{"key": "value"}
		var expected bytes.Buffer
		err := json.NewEncoder(&expected).Encode(data)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		var buf bytes.Buffer
		err = jsonProtocol.Serialize(&buf, data)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !reflect.DeepEqual(buf.Bytes(), expected.Bytes()) {
			t.Errorf("Expected %v, got %v", expected, buf.Bytes())
		}
	})

	t.Run("Deserialize", func(t *testing.T) {
		jsonData := []byte(`{"key":"value"}`)
		var result map[string]string

		err := jsonProtocol.Deserialize(bytes.NewReader(jsonData), &result)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := map[string]string{"key": "value"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}
