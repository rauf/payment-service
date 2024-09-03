package registry

import (
	"errors"
	"sync"
)

// Registry is a generic registry that stores values by name and also stores the order of registration.
type Registry[T any] struct {
	registry map[string]T
	order    []string
	mu       sync.RWMutex
}

func NewRegistry[T any]() *Registry[T] {
	return &Registry[T]{
		registry: make(map[string]T),
		order:    []string{},
	}
}

func (r *Registry[T]) Register(name string, value T) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.registry[name]; exists {
		return errors.New("value already registered")
	}

	r.registry[name] = value
	r.order = append(r.order, name)
	return nil
}

func (r *Registry[T]) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.registry[name]; !exists {
		return errors.New("value not found")
	}

	delete(r.registry, name)
	for i, n := range r.order {
		if n == name {
			r.order = append(r.order[:i], r.order[i+1:]...)
			break
		}
	}
	return nil
}

func (r *Registry[T]) Get(name string) (T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var zero T
	g, exists := r.registry[name]
	if !exists {
		return zero, errors.New("value not found")
	}
	return g, nil
}

func (r *Registry[T]) List() []T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	values := make([]T, 0, len(r.registry))
	for _, o := range r.order {
		value, err := r.Get(o)
		if err != nil {
			continue
		}
		values = append(values, value)
	}
	return values
}

func (r *Registry[T]) SetOrder(order []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(order) != len(r.order) {
		return errors.New("order length does not match registry length")
	}

	for _, name := range order {
		if _, exists := r.registry[name]; !exists {
			return errors.New("invalid element in order: " + name)
		}
	}

	r.order = order
	return nil
}

// ListWithPreference returns the list with the preferred at the beginning.
func (r *Registry[T]) ListWithPreference(preferred string) ([]T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.order) == 0 {
		return nil, errors.New("no registry registered")
	}

	if preferred == "" {
		return r.List(), nil
	}

	preferredIndex := -1
	for i, name := range r.order {
		if name == preferred {
			preferredIndex = i
			break
		}
	}

	if preferredIndex == -1 {
		return r.List(), nil
	}

	ordered := make([]T, 0, len(r.registry))
	ordered = append(ordered, r.registry[preferred])
	for i := 0; i < len(r.order); i++ {
		if i == preferredIndex {
			continue
		}
		ordered = append(ordered, r.registry[r.order[i]])
	}
	return ordered, nil
}
