package core

import "fmt"

type Registry struct {
	services map[string]ServiceDefinition
}

func NewRegistry() *Registry {
	return &Registry{services: map[string]ServiceDefinition{}}
}

func (r *Registry) RegisterService(service ServiceDefinition) error {
	if service.ID == "" {
		return fmt.Errorf("service id is required")
	}
	if _, exists := r.services[service.ID]; exists {
		return fmt.Errorf("service %q already registered", service.ID)
	}
	r.services[service.ID] = service
	return nil
}

func (r *Registry) Service(id string) (ServiceDefinition, bool) {
	service, ok := r.services[id]
	return service, ok
}

func (r *Registry) Operation(serviceID string, operationName string) (OperationDefinition, bool) {
	service, ok := r.Service(serviceID)
	if !ok {
		return OperationDefinition{}, false
	}
	for _, operation := range service.Operations {
		if operation.Name == operationName {
			return operation, true
		}
	}
	return OperationDefinition{}, false
}

func (r *Registry) Operations() []RegisteredOperation {
	operations := []RegisteredOperation{}
	for _, service := range r.services {
		for _, operation := range service.Operations {
			operations = append(operations, RegisteredOperation{
				Service:   service,
				Operation: operation,
			})
		}
	}
	return operations
}
