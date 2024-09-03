package registry

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

func TestNewRegistry(t *testing.T) {
	r := NewRegistry[string]()
	if r == nil {
		t.Fatal("NewRegistry returned nil")
	}
	if len(r.registry) != 0 {
		t.Errorf("Expected empty registry, got %d items", len(r.registry))
	}
	if len(r.order) != 0 {
		t.Errorf("Expected empty order, got %d items", len(r.order))
	}
}

func TestRegister(t *testing.T) {
	r := NewRegistry[string]()

	tests := []struct {
		name        string
		key         string
		value       string
		expectError bool
	}{
		{"Register first item", "key1", "value1", false},
		{"Register second item", "key2", "value2", false},
		{"Register duplicate key", "key1", "value3", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := r.Register(tt.key, tt.value)
			if (err != nil) != tt.expectError {
				t.Errorf("Register() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}

	if len(r.registry) != 2 {
		t.Errorf("Expected 2 items in registry, got %d", len(r.registry))
	}
	if len(r.order) != 2 {
		t.Errorf("Expected 2 items in order, got %d", len(r.order))
	}
}

func TestUnregister(t *testing.T) {
	r := NewRegistry[string]()
	_ = r.Register("key1", "value1")
	_ = r.Register("key2", "value2")

	tests := []struct {
		name        string
		key         string
		expectError bool
	}{
		{"Unregister existing item", "key1", false},
		{"Unregister non-existing item", "key3", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := r.Unregister(tt.key)
			if (err != nil) != tt.expectError {
				t.Errorf("Unregister() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}

	if len(r.registry) != 1 {
		t.Errorf("Expected 1 item in registry, got %d", len(r.registry))
	}
	if len(r.order) != 1 {
		t.Errorf("Expected 1 item in order, got %d", len(r.order))
	}
}

func TestGet(t *testing.T) {
	r := NewRegistry[string]()
	_ = r.Register("key1", "value1")

	tests := []struct {
		name        string
		key         string
		expectValue string
		expectError bool
	}{
		{"Get existing item", "key1", "value1", false},
		{"Get non-existing item", "key2", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := r.Get(tt.key)
			if (err != nil) != tt.expectError {
				t.Errorf("Get() error = %v, expectError %v", err, tt.expectError)
			}
			if value != tt.expectValue {
				t.Errorf("Get() value = %v, expected %v", value, tt.expectValue)
			}
		})
	}
}

func TestList(t *testing.T) {
	r := NewRegistry[string]()
	_ = r.Register("key1", "value1")
	_ = r.Register("key2", "value2")
	_ = r.Register("key3", "value3")

	list := r.List()
	expected := []string{"value1", "value2", "value3"}

	if !reflect.DeepEqual(list, expected) {
		t.Errorf("List() = %v, expected %v", list, expected)
	}
}

func TestSetOrder(t *testing.T) {
	r := NewRegistry[string]()
	_ = r.Register("key1", "value1")
	_ = r.Register("key2", "value2")
	_ = r.Register("key3", "value3")

	tests := []struct {
		name        string
		order       []string
		expectError bool
	}{
		{"Valid order", []string{"key2", "key3", "key1"}, false},
		{"Invalid order (missing key)", []string{"key2", "key3"}, true},
		{"Invalid order (extra key)", []string{"key1", "key2", "key3", "key4"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := r.SetOrder(tt.order)
			if (err != nil) != tt.expectError {
				t.Errorf("SetOrder() error = %v, expectError %v", err, tt.expectError)
			}
			if err == nil && !reflect.DeepEqual(r.order, tt.order) {
				t.Errorf("SetOrder() order = %v, expected %v", r.order, tt.order)
			}
		})
	}
}

func TestListWithPreference(t *testing.T) {
	r := NewRegistry[string]()
	_ = r.Register("key1", "value1")
	_ = r.Register("key2", "value2")
	_ = r.Register("key3", "value3")

	tests := []struct {
		name        string
		preferred   string
		expectOrder []string
		expectError bool
	}{
		{"Prefer first item", "key1", []string{"value1", "value2", "value3"}, false},
		{"Prefer middle item", "key2", []string{"value2", "value1", "value3"}, false},
		{"Prefer last item", "key3", []string{"value3", "value1", "value2"}, false},
		{"Prefer non-existing item", "key4", []string{"value1", "value2", "value3"}, false},
		{"No preference", "", []string{"value1", "value2", "value3"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := r.ListWithPreference(tt.preferred)
			if (err != nil) != tt.expectError {
				t.Errorf("ListWithPreference() error = %v, expectError %v", err, tt.expectError)
			}
			if !reflect.DeepEqual(list, tt.expectOrder) {
				t.Errorf("ListWithPreference() = %v, expected %v", list, tt.expectOrder)
			}
		})
	}
}

func TestConcurrency(t *testing.T) {
	r := NewRegistry[int]()
	const goroutines = 100
	const operations = 1000

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				key := fmt.Sprintf("key%d-%d", id, j)
				_ = r.Register(key, j)
				_, _ = r.Get(key)
				_ = r.Unregister(key)
			}
		}(i)
	}

	wg.Wait()

	if len(r.registry) != 0 {
		t.Errorf("Expected empty registry after concurrent operations, got %d items", len(r.registry))
	}
}

func TestRegistryWithCustomType(t *testing.T) {
	type CustomType struct {
		ID   int
		Name string
	}

	r := NewRegistry[CustomType]()

	item1 := CustomType{ID: 1, Name: "Item 1"}
	item2 := CustomType{ID: 2, Name: "Item 2"}

	_ = r.Register("item1", item1)
	_ = r.Register("item2", item2)

	retrieved, err := r.Get("item1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !reflect.DeepEqual(retrieved, item1) {
		t.Errorf("Retrieved item doesn't match: got %v, want %v", retrieved, item1)
	}

	list := r.List()
	expected := []CustomType{item1, item2}
	if !reflect.DeepEqual(list, expected) {
		t.Errorf("List() = %v, expected %v", list, expected)
	}
}
