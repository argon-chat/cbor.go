package cbor

import (
	"errors"
	"fmt"
)

// Common CBOR errors.
var (
	// ErrUnexpectedEndOfData is returned when the data ends unexpectedly.
	ErrUnexpectedEndOfData = errors.New("cbor: unexpected end of data")

	// ErrInvalidCbor is returned when the CBOR data is malformed.
	ErrInvalidCbor = errors.New("cbor: invalid CBOR data")

	// ErrInvalidMajorType is returned when an unexpected major type is encountered.
	ErrInvalidMajorType = errors.New("cbor: invalid major type")

	// ErrInvalidSimpleValue is returned when an invalid simple value is encountered.
	ErrInvalidSimpleValue = errors.New("cbor: invalid simple value")

	// ErrInvalidUtf8 is returned when a text string contains invalid UTF-8.
	ErrInvalidUtf8 = errors.New("cbor: invalid UTF-8 in text string")

	// ErrOverflow is returned when a value overflows the target type.
	ErrOverflow = errors.New("cbor: integer overflow")

	// ErrUnexpectedBreak is returned when a break byte is encountered unexpectedly.
	ErrUnexpectedBreak = errors.New("cbor: unexpected break")

	// ErrNonCanonical is returned in strict/canonical mode when encoding is non-canonical.
	ErrNonCanonical = errors.New("cbor: non-canonical encoding")

	// ErrNotAtEnd is returned when there is remaining data after the root value.
	ErrNotAtEnd = errors.New("cbor: unexpected data after root value")

	// ErrInvalidState is returned when an operation is attempted in an invalid state.
	ErrInvalidState = errors.New("cbor: invalid reader state for this operation")

	// ErrDuplicateKey is returned when a duplicate key is found in a map (in strict mode).
	ErrDuplicateKey = errors.New("cbor: duplicate key in map")

	// ErrUnsortedKeys is returned when map keys are not sorted (in canonical mode).
	ErrUnsortedKeys = errors.New("cbor: map keys are not sorted")

	// ErrIndefiniteLengthNotAllowed is returned when indefinite length is used in canonical mode.
	ErrIndefiniteLengthNotAllowed = errors.New("cbor: indefinite length not allowed in canonical mode")

	// ErrBufferTooSmall is returned when the buffer is too small for the operation.
	ErrBufferTooSmall = errors.New("cbor: buffer too small")

	// ErrNestingDepthExceeded is returned when the maximum nesting depth is exceeded.
	ErrNestingDepthExceeded = errors.New("cbor: maximum nesting depth exceeded")

	// ErrMissingBreak is returned when an indefinite-length item is not terminated.
	ErrMissingBreak = errors.New("cbor: missing break for indefinite-length item")

	// ErrIncompleteContainer is returned when a container has fewer items than expected.
	ErrIncompleteContainer = errors.New("cbor: incomplete container")

	// ErrExtraItems is returned when a container has more items than expected.
	ErrExtraItems = errors.New("cbor: extra items in container")
)

// CborError provides detailed error information.
type CborError struct {
	Err     error
	Offset  int
	Message string
}

// Error implements the error interface.
func (e *CborError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("cbor error at offset %d: %s: %v", e.Offset, e.Message, e.Err)
	}
	return fmt.Sprintf("cbor error at offset %d: %v", e.Offset, e.Err)
}

// Unwrap returns the underlying error.
func (e *CborError) Unwrap() error {
	return e.Err
}

// NewCborError creates a new CborError.
func NewCborError(err error, offset int, message string) *CborError {
	return &CborError{
		Err:     err,
		Offset:  offset,
		Message: message,
	}
}

// TypeMismatchError is returned when the expected type doesn't match the actual type.
type TypeMismatchError struct {
	Expected CborReaderState
	Actual   CborReaderState
}

// Error implements the error interface.
func (e *TypeMismatchError) Error() string {
	return fmt.Sprintf("cbor: expected %s but got %s", e.Expected, e.Actual)
}
