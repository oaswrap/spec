package spec

import "sync"

// SpecError is a thread-safe error collector for OpenAPI specifications.
type SpecError struct {
	mu     sync.Mutex
	errors []error
}

func (se *SpecError) add(err error) {
	se.mu.Lock()
	defer se.mu.Unlock()
	if err != nil {
		se.errors = append(se.errors, err)
	}
}

// Errors returns a slice of collected errors.
func (se *SpecError) Errors() []error {
	se.mu.Lock()
	defer se.mu.Unlock()
	return se.errors
}

// Error implements the error interface for SpecError.
func (se *SpecError) Error() string {
	se.mu.Lock()
	defer se.mu.Unlock()
	if len(se.errors) == 0 {
		return ""
	}
	result := "Spec errors:\n"
	for _, err := range se.errors {
		result += "- " + err.Error() + "\n"
	}
	return result
}

// HasErrors checks if there are any collected errors.
func (se *SpecError) HasErrors() bool {
	se.mu.Lock()
	defer se.mu.Unlock()
	return len(se.errors) > 0
}
