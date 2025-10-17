package errs_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/oaswrap/spec/internal/errs"
	"github.com/stretchr/testify/assert"
)

func TestSpecError_Add(t *testing.T) {
	se := &errs.SpecError{}

	// Test adding a valid error
	err := errors.New("test error")
	se.Add(err)

	assert.Len(t, se.Errors(), 1)

	// Test adding nil error (should not be added)
	se.Add(nil)

	assert.Len(t, se.Errors(), 1)
}

func TestSpecError_Errors(t *testing.T) {
	se := &errs.SpecError{}

	// Test empty errors
	errs := se.Errors()
	assert.Empty(t, errs)

	// Test with errors
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")
	se.Add(err1)
	se.Add(err2)

	errs = se.Errors()
	assert.Len(t, errs, 2)
}

func TestSpecError_Error(t *testing.T) {
	se := &errs.SpecError{}

	// Test empty error message
	msg := se.Error()
	assert.Empty(t, msg)

	// Test with single error
	se.Add(errors.New("test error"))
	msg = se.Error()
	expected := "Spec errors:\n- test error\n"
	assert.Equal(t, expected, msg)

	// Test with multiple errors
	se.Add(errors.New("second error"))
	msg = se.Error()
	expected = "Spec errors:\n- test error\n- second error\n"
	assert.Equal(t, expected, msg)
}

func TestSpecError_HasErrors(t *testing.T) {
	se := &errs.SpecError{}

	// Test with no errors
	assert.False(t, se.HasErrors())

	// Test with errors
	se.Add(errors.New("test error"))
	assert.True(t, se.HasErrors())
}

func TestSpecError_ConcurrentAccess(t *testing.T) {
	se := &errs.SpecError{}
	var wg sync.WaitGroup

	// Test concurrent adds
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			se.Add(errors.New("error"))
		}()
	}

	wg.Wait()

	assert.Len(t, se.Errors(), 100)
}
