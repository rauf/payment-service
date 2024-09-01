package gateway

import (
	"errors"
	"sync"
)

type Registry struct {
	gateways map[string]PaymentGateway
	order    []string
	mu       sync.RWMutex
}

func NewGatewayRegistry() *Registry {
	return &Registry{
		gateways: make(map[string]PaymentGateway),
		order:    []string{},
	}
}

func (r *Registry) Register(name string, gateway PaymentGateway) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.gateways[name]; exists {
		return errors.New("gateway already registered")
	}

	r.gateways[name] = gateway
	r.order = append(r.order, name)
	return nil
}

func (r *Registry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.gateways[name]; !exists {
		return errors.New("gateway not found")
	}

	delete(r.gateways, name)
	for i, n := range r.order {
		if n == name {
			r.order = append(r.order[:i], r.order[i+1:]...)
			break
		}
	}
	return nil
}

func (r *Registry) Get(name string) (PaymentGateway, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	g, exists := r.gateways[name]
	if !exists {
		return nil, errors.New("gateway not found")
	}
	return g, nil
}

func (r *Registry) List() []PaymentGateway {
	r.mu.RLock()
	defer r.mu.RUnlock()

	gateways := make([]PaymentGateway, 0, len(r.gateways))
	for _, o := range r.order {
		gateway, err := r.Get(o)
		if err != nil {
			continue
		}
		gateways = append(gateways, gateway)
	}
	return gateways
}

func (r *Registry) SetOrder(order []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, name := range order {
		if _, exists := r.gateways[name]; !exists {
			return errors.New("invalid gateway in order: " + name)
		}
	}

	r.order = order
	return nil
}

// ListWithPreference returns the list of gateways with the preferred gateway at the beginning.
func (r *Registry) ListWithPreference(preferred string) ([]PaymentGateway, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.order) == 0 {
		return nil, errors.New("no gateways registered")
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
		return nil, errors.New("preferred gateway not found")
	}

	ordered := make([]PaymentGateway, 0, len(r.gateways))
	ordered = append(ordered, r.gateways[preferred])
	for i := 0; i < len(r.order); i++ {
		if i == preferredIndex {
			continue
		}
		ordered = append(ordered, r.gateways[r.order[i]])
	}
	return ordered, nil
}
