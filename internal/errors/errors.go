package errors

import "sync"

type MultiError struct {
	mu     sync.Mutex
	errors []error
}

func (e *MultiError) Add(err error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.errors = append(e.errors, err)
}

func (e *MultiError) Errors() []error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.errors
}

func (e *MultiError) Error() string {
	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.errors) == 0 {
		return ""
	}

	var errStr string
	for _, err := range e.errors {
		errStr += err.Error() + "\n"
	}
	return errStr
}

func (e *MultiError) HasErrors() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.errors) > 0
}
