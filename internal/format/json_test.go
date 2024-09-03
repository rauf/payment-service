package format

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestJSONProtocol(t *testing.T) {
	jsonProtocol := NewJSONProtocol()

	t.Run("Marshal", func(t *testing.T) {
		data := map[string]string{"key": "value"}
		expected, _ := json.Marshal(data)

		result, err := jsonProtocol.Marshal(data)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Unmarshal", func(t *testing.T) {
		jsonData := []byte(`{"key":"value"}`)
		var result map[string]string

		err := jsonProtocol.Unmarshal(jsonData, &result)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expected := map[string]string{"key": "value"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}
