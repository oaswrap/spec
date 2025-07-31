package spec

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpecError_add(t *testing.T) {
	se := &SpecError{}

	// Test adding a valid error
	err := errors.New("test error")
	se.add(err)

	assert.Len(t, se.errors, 1)

	// Test adding nil error (should not be added)
	se.add(nil)

	assert.Len(t, se.errors, 1)
}

func TestSpecError_Errors(t *testing.T) {
	se := &SpecError{}

	// Test empty errors
	errs := se.Errors()
	assert.Empty(t, errs)

	// Test with errors
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")
	se.add(err1)
	se.add(err2)

	errs = se.Errors()
	assert.Len(t, errs, 2)
}

func TestSpecError_Error(t *testing.T) {
	se := &SpecError{}

	// Test empty error message
	msg := se.Error()
	assert.Empty(t, msg)

	// Test with single error
	se.add(errors.New("test error"))
	msg = se.Error()
	expected := "Spec errors:\n- test error\n"
	assert.Equal(t, expected, msg)

	// Test with multiple errors
	se.add(errors.New("second error"))
	msg = se.Error()
	expected = "Spec errors:\n- test error\n- second error\n"
	assert.Equal(t, expected, msg)
}

func TestSpecError_HasErrors(t *testing.T) {
	se := &SpecError{}

	// Test with no errors
	assert.False(t, se.HasErrors())

	// Test with errors
	se.add(errors.New("test error"))
	assert.True(t, se.HasErrors())
}

func TestSpecError_ConcurrentAccess(t *testing.T) {
	se := &SpecError{}
	var wg sync.WaitGroup

	// Test concurrent adds
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			se.add(errors.New("error"))
		}(i)
	}

	wg.Wait()

	assert.Len(t, se.Errors(), 100)
}
