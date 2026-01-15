package cbor

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/big"
	"time"
	"unicode/utf8"
)

// CborReader provides methods for reading CBOR encoded data.
type CborReader struct {
	data                    []byte
	offset                  int
	conformanceMode         CborConformanceMode
	nestingStack            []readerNestingInfo
	maxNestingDepth         int
	cachedState             CborReaderState
	stateComputed           bool
	allowMultipleRootValues bool
}

// readerNestingInfo tracks the state of nested containers during reading.
type readerNestingInfo struct {
	majorType      MajorType
	definiteLength int64 // -1 for indefinite
	itemsRead      int64
	isMap          bool
	keyRead        bool // for maps, tracks if we're expecting a value
	isIndefinite   bool
}

// ReaderOption is a function that configures a CborReader.
type ReaderOption func(*CborReader)

// WithReaderConformanceMode sets the conformance mode for the reader.
func WithReaderConformanceMode(mode CborConformanceMode) ReaderOption {
	return func(r *CborReader) {
		r.conformanceMode = mode
	}
}

// WithReaderMaxNestingDepth sets the maximum nesting depth for the reader.
func WithReaderMaxNestingDepth(depth int) ReaderOption {
	return func(r *CborReader) {
		r.maxNestingDepth = depth
	}
}

// WithReaderAllowMultipleRootValues allows reading multiple root-level values.
func WithReaderAllowMultipleRootValues(allow bool) ReaderOption {
	return func(r *CborReader) {
		r.allowMultipleRootValues = allow
	}
}

// NewCborReader creates a new CborReader for the given data.
func NewCborReader(data []byte, opts ...ReaderOption) *CborReader {
	r := &CborReader{
		data:            data,
		offset:          0,
		conformanceMode: ConformanceLax,
		nestingStack:    make([]readerNestingInfo, 0, 16),
		maxNestingDepth: 64,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

// Reset resets the reader to the beginning.
func (r *CborReader) Reset() {
	r.offset = 0
	r.nestingStack = r.nestingStack[:0]
	r.cachedState = StateUndefined
	r.stateComputed = false
}

// ResetWithData resets the reader with new data.
func (r *CborReader) ResetWithData(data []byte) {
	r.data = data
	r.Reset()
}

// BytesRemaining returns the number of bytes remaining to be read.
func (r *CborReader) BytesRemaining() int {
	return len(r.data) - r.offset
}

// CurrentOffset returns the current position in the data.
func (r *CborReader) CurrentOffset() int {
	return r.offset
}

// NestingDepth returns the current nesting depth.
func (r *CborReader) NestingDepth() int {
	return len(r.nestingStack)
}

// invalidateState clears the cached state.
func (r *CborReader) invalidateState() {
	r.stateComputed = false
}

// PeekState returns the current state without advancing the reader.
func (r *CborReader) PeekState() (CborReaderState, error) {
	if r.stateComputed {
		return r.cachedState, nil
	}

	state, err := r.computeState()
	if err != nil {
		return StateUndefined, err
	}

	r.cachedState = state
	r.stateComputed = true
	return state, nil
}

// computeState determines the current reader state.
func (r *CborReader) computeState() (CborReaderState, error) {
	// Check if we're at the end of a container
	if len(r.nestingStack) > 0 {
		info := &r.nestingStack[len(r.nestingStack)-1]

		if !info.isIndefinite && info.itemsRead >= info.definiteLength {
			if info.isMap {
				return StateEndMap, nil
			}
			return StateEndArray, nil
		}
	}

	if r.offset >= len(r.data) {
		if len(r.nestingStack) > 0 {
			return StateUndefined, ErrUnexpectedEndOfData
		}
		return StateFinished, nil
	}

	initialByte := r.data[r.offset]

	// Check for break byte
	if initialByte == breakByte {
		if len(r.nestingStack) == 0 {
			return StateUndefined, ErrUnexpectedBreak
		}

		info := &r.nestingStack[len(r.nestingStack)-1]
		if !info.isIndefinite {
			return StateUndefined, ErrUnexpectedBreak
		}

		switch info.majorType {
		case MajorTypeArray:
			return StateEndArray, nil
		case MajorTypeMap:
			if info.keyRead {
				return StateUndefined, ErrIncompleteContainer
			}
			return StateEndMap, nil
		case MajorTypeByteString:
			return StateEndIndefiniteLengthByteString, nil
		case MajorTypeTextString:
			return StateEndIndefiniteLengthTextString, nil
		}
	}

	mt, ai := decodeInitialByte(initialByte)

	switch mt {
	case MajorTypeUnsignedInteger:
		return StateUnsignedInteger, nil
	case MajorTypeNegativeInteger:
		return StateNegativeInteger, nil
	case MajorTypeByteString:
		if ai == byte(AdditionalInfoIndefiniteLength) {
			return StateStartIndefiniteLengthByteString, nil
		}
		return StateByteString, nil
	case MajorTypeTextString:
		if ai == byte(AdditionalInfoIndefiniteLength) {
			return StateStartIndefiniteLengthTextString, nil
		}
		return StateTextString, nil
	case MajorTypeArray:
		return StateStartArray, nil
	case MajorTypeMap:
		return StateStartMap, nil
	case MajorTypeTag:
		return StateTag, nil
	case MajorTypeSimpleOrFloat:
		switch ai {
		case byte(SimpleValueFalse), byte(SimpleValueTrue):
			return StateBoolean, nil
		case byte(SimpleValueNull):
			return StateNull, nil
		case byte(SimpleValueUndefined):
			return StateUndefinedValue, nil
		case 24:
			return StateSimpleValue, nil
		case 25:
			return StateHalfPrecisionFloat, nil
		case 26:
			return StateSinglePrecisionFloat, nil
		case 27:
			return StateDoublePrecisionFloat, nil
		default:
			if ai < 24 {
				return StateSimpleValue, nil
			}
			return StateUndefined, ErrInvalidSimpleValue
		}
	}

	return StateUndefined, ErrInvalidMajorType
}

// readInitialByte reads the initial byte and returns the additional information value.
func (r *CborReader) readArgumentValue(mt MajorType) (uint64, error) {
	if r.offset >= len(r.data) {
		return 0, ErrUnexpectedEndOfData
	}

	initialByte := r.data[r.offset]
	actualMt, ai := decodeInitialByte(initialByte)

	if actualMt != mt {
		return 0, &TypeMismatchError{Expected: CborReaderState(mt), Actual: CborReaderState(actualMt)}
	}

	r.offset++

	switch {
	case ai < 24:
		return uint64(ai), nil
	case ai == 24:
		if r.offset >= len(r.data) {
			return 0, ErrUnexpectedEndOfData
		}
		val := r.data[r.offset]
		r.offset++

		// Canonical check: value must be >= 24
		if r.conformanceMode >= ConformanceStrict && val < 24 {
			return 0, ErrNonCanonical
		}
		return uint64(val), nil
	case ai == 25:
		if r.offset+2 > len(r.data) {
			return 0, ErrUnexpectedEndOfData
		}
		val := binary.BigEndian.Uint16(r.data[r.offset:])
		r.offset += 2

		// Canonical check: value must be > 255
		if r.conformanceMode >= ConformanceStrict && val <= 0xFF {
			return 0, ErrNonCanonical
		}
		return uint64(val), nil
	case ai == 26:
		if r.offset+4 > len(r.data) {
			return 0, ErrUnexpectedEndOfData
		}
		val := binary.BigEndian.Uint32(r.data[r.offset:])
		r.offset += 4

		// Canonical check: value must be > 65535
		if r.conformanceMode >= ConformanceStrict && val <= 0xFFFF {
			return 0, ErrNonCanonical
		}
		return uint64(val), nil
	case ai == 27:
		if r.offset+8 > len(r.data) {
			return 0, ErrUnexpectedEndOfData
		}
		val := binary.BigEndian.Uint64(r.data[r.offset:])
		r.offset += 8

		// Canonical check: value must be > 4294967295
		if r.conformanceMode >= ConformanceStrict && val <= 0xFFFFFFFF {
			return 0, ErrNonCanonical
		}
		return uint64(val), nil
	case ai == 31:
		return 0, nil // Indefinite length
	default:
		return 0, ErrInvalidCbor
	}
}

// advanceContainer updates container state after reading an item.
func (r *CborReader) advanceContainer() {
	if len(r.nestingStack) == 0 {
		return
	}

	info := &r.nestingStack[len(r.nestingStack)-1]
	if info.isMap {
		if info.keyRead {
			// We just read a value
			info.keyRead = false
			info.itemsRead++
		} else {
			// We just read a key
			info.keyRead = true
		}
	} else {
		info.itemsRead++
	}
	r.invalidateState()
}

// ReadUint64 reads an unsigned 64-bit integer.
func (r *CborReader) ReadUint64() (uint64, error) {
	state, err := r.PeekState()
	if err != nil {
		return 0, err
	}
	if state != StateUnsignedInteger {
		return 0, &TypeMismatchError{Expected: StateUnsignedInteger, Actual: state}
	}

	r.invalidateState()
	val, err := r.readArgumentValue(MajorTypeUnsignedInteger)
	if err != nil {
		return 0, err
	}

	r.advanceContainer()
	return val, nil
}

// ReadInt64 reads a signed 64-bit integer (can be positive or negative).
func (r *CborReader) ReadInt64() (int64, error) {
	state, err := r.PeekState()
	if err != nil {
		return 0, err
	}

	r.invalidateState()

	switch state {
	case StateUnsignedInteger:
		val, err := r.readArgumentValue(MajorTypeUnsignedInteger)
		if err != nil {
			return 0, err
		}
		if val > math.MaxInt64 {
			return 0, ErrOverflow
		}
		r.advanceContainer()
		return int64(val), nil

	case StateNegativeInteger:
		val, err := r.readArgumentValue(MajorTypeNegativeInteger)
		if err != nil {
			return 0, err
		}
		// CBOR negative integers are encoded as -1 - n
		if val > math.MaxInt64 {
			return 0, ErrOverflow
		}
		r.advanceContainer()
		return -1 - int64(val), nil

	default:
		return 0, &TypeMismatchError{Expected: StateUnsignedInteger, Actual: state}
	}
}

// ReadInt32 reads a signed 32-bit integer.
func (r *CborReader) ReadInt32() (int32, error) {
	val, err := r.ReadInt64()
	if err != nil {
		return 0, err
	}
	if val < math.MinInt32 || val > math.MaxInt32 {
		return 0, ErrOverflow
	}
	return int32(val), nil
}

// ReadUint32 reads an unsigned 32-bit integer.
func (r *CborReader) ReadUint32() (uint32, error) {
	val, err := r.ReadUint64()
	if err != nil {
		return 0, err
	}
	if val > math.MaxUint32 {
		return 0, ErrOverflow
	}
	return uint32(val), nil
}

// ReadInt16 reads a signed 16-bit integer.
func (r *CborReader) ReadInt16() (int16, error) {
	val, err := r.ReadInt64()
	if err != nil {
		return 0, err
	}
	if val < math.MinInt16 || val > math.MaxInt16 {
		return 0, ErrOverflow
	}
	return int16(val), nil
}

// ReadUint16 reads an unsigned 16-bit integer.
func (r *CborReader) ReadUint16() (uint16, error) {
	val, err := r.ReadUint64()
	if err != nil {
		return 0, err
	}
	if val > math.MaxUint16 {
		return 0, ErrOverflow
	}
	return uint16(val), nil
}

// ReadInt8 reads a signed 8-bit integer.
func (r *CborReader) ReadInt8() (int8, error) {
	val, err := r.ReadInt64()
	if err != nil {
		return 0, err
	}
	if val < math.MinInt8 || val > math.MaxInt8 {
		return 0, ErrOverflow
	}
	return int8(val), nil
}

// ReadUint8 reads an unsigned 8-bit integer.
func (r *CborReader) ReadUint8() (uint8, error) {
	val, err := r.ReadUint64()
	if err != nil {
		return 0, err
	}
	if val > math.MaxUint8 {
		return 0, ErrOverflow
	}
	return uint8(val), nil
}

// ReadInt reads an int value.
func (r *CborReader) ReadInt() (int, error) {
	val, err := r.ReadInt64()
	if err != nil {
		return 0, err
	}
	// Check for overflow on 32-bit systems
	if val < math.MinInt || val > math.MaxInt {
		return 0, ErrOverflow
	}
	return int(val), nil
}

// ReadBigInt reads an integer as a big.Int, handling bignums if tagged.
func (r *CborReader) ReadBigInt() (*big.Int, error) {
	state, err := r.PeekState()
	if err != nil {
		return nil, err
	}

	switch state {
	case StateUnsignedInteger:
		val, err := r.ReadUint64()
		if err != nil {
			return nil, err
		}
		return new(big.Int).SetUint64(val), nil

	case StateNegativeInteger:
		val, err := r.ReadInt64()
		if err != nil {
			// Might be a bignum that overflows int64
			r.stateComputed = true
			r.cachedState = StateNegativeInteger
			// Read as raw negative value
			r.invalidateState()
			raw, err2 := r.readArgumentValue(MajorTypeNegativeInteger)
			if err2 != nil {
				return nil, err2
			}
			r.advanceContainer()
			// -1 - raw
			result := new(big.Int).SetUint64(raw)
			result.Add(result, big.NewInt(1))
			result.Neg(result)
			return result, nil
		}
		return big.NewInt(val), nil

	case StateTag:
		tag, err := r.ReadTag()
		if err != nil {
			return nil, err
		}

		switch tag {
		case TagUnsignedBignum:
			data, err := r.ReadByteString()
			if err != nil {
				return nil, err
			}
			return new(big.Int).SetBytes(data), nil

		case TagNegativeBignum:
			data, err := r.ReadByteString()
			if err != nil {
				return nil, err
			}
			// -1 - n
			result := new(big.Int).SetBytes(data)
			result.Add(result, big.NewInt(1))
			result.Neg(result)
			return result, nil

		default:
			return nil, &TypeMismatchError{Expected: StateUnsignedInteger, Actual: StateTag}
		}

	default:
		return nil, &TypeMismatchError{Expected: StateUnsignedInteger, Actual: state}
	}
}

// ReadByteString reads a byte string.
func (r *CborReader) ReadByteString() ([]byte, error) {
	state, err := r.PeekState()
	if err != nil {
		return nil, err
	}

	if state == StateStartIndefiniteLengthByteString {
		return r.readIndefiniteByteString()
	}

	if state != StateByteString {
		return nil, &TypeMismatchError{Expected: StateByteString, Actual: state}
	}

	r.invalidateState()
	length, err := r.readArgumentValue(MajorTypeByteString)
	if err != nil {
		return nil, err
	}

	if r.offset+int(length) > len(r.data) {
		return nil, ErrUnexpectedEndOfData
	}

	result := make([]byte, length)
	copy(result, r.data[r.offset:r.offset+int(length)])
	r.offset += int(length)
	r.advanceContainer()
	return result, nil
}

// readIndefiniteByteString reads an indefinite-length byte string.
func (r *CborReader) readIndefiniteByteString() ([]byte, error) {
	if r.conformanceMode >= ConformanceCanonical {
		return nil, ErrIndefiniteLengthNotAllowed
	}

	// Skip the initial byte
	r.offset++
	r.invalidateState()

	var result bytes.Buffer

	for {
		if r.offset >= len(r.data) {
			return nil, ErrUnexpectedEndOfData
		}

		if r.data[r.offset] == breakByte {
			r.offset++
			break
		}

		// Read a definite-length byte string chunk
		mt, _ := decodeInitialByte(r.data[r.offset])
		if mt != MajorTypeByteString {
			return nil, ErrInvalidCbor
		}

		length, err := r.readArgumentValue(MajorTypeByteString)
		if err != nil {
			return nil, err
		}

		if r.offset+int(length) > len(r.data) {
			return nil, ErrUnexpectedEndOfData
		}

		result.Write(r.data[r.offset : r.offset+int(length)])
		r.offset += int(length)
	}

	r.advanceContainer()
	return result.Bytes(), nil
}

// ReadTextString reads a UTF-8 text string.
func (r *CborReader) ReadTextString() (string, error) {
	state, err := r.PeekState()
	if err != nil {
		return "", err
	}

	if state == StateStartIndefiniteLengthTextString {
		return r.readIndefiniteTextString()
	}

	if state != StateTextString {
		return "", &TypeMismatchError{Expected: StateTextString, Actual: state}
	}

	r.invalidateState()
	length, err := r.readArgumentValue(MajorTypeTextString)
	if err != nil {
		return "", err
	}

	if r.offset+int(length) > len(r.data) {
		return "", ErrUnexpectedEndOfData
	}

	strBytes := r.data[r.offset : r.offset+int(length)]

	// Validate UTF-8 in strict mode
	if r.conformanceMode >= ConformanceStrict && !utf8.Valid(strBytes) {
		return "", ErrInvalidUtf8
	}

	result := string(strBytes)
	r.offset += int(length)
	r.advanceContainer()
	return result, nil
}

// readIndefiniteTextString reads an indefinite-length text string.
func (r *CborReader) readIndefiniteTextString() (string, error) {
	if r.conformanceMode >= ConformanceCanonical {
		return "", ErrIndefiniteLengthNotAllowed
	}

	// Skip the initial byte
	r.offset++
	r.invalidateState()

	var result bytes.Buffer

	for {
		if r.offset >= len(r.data) {
			return "", ErrUnexpectedEndOfData
		}

		if r.data[r.offset] == breakByte {
			r.offset++
			break
		}

		// Read a definite-length text string chunk
		mt, _ := decodeInitialByte(r.data[r.offset])
		if mt != MajorTypeTextString {
			return "", ErrInvalidCbor
		}

		length, err := r.readArgumentValue(MajorTypeTextString)
		if err != nil {
			return "", err
		}

		if r.offset+int(length) > len(r.data) {
			return "", ErrUnexpectedEndOfData
		}

		chunk := r.data[r.offset : r.offset+int(length)]

		if r.conformanceMode >= ConformanceStrict && !utf8.Valid(chunk) {
			return "", ErrInvalidUtf8
		}

		result.Write(chunk)
		r.offset += int(length)
	}

	r.advanceContainer()
	return result.String(), nil
}

// ReadStartArray reads the start of an array and returns its length.
// Returns -1 for indefinite-length arrays.
func (r *CborReader) ReadStartArray() (int, error) {
	state, err := r.PeekState()
	if err != nil {
		return 0, err
	}
	if state != StateStartArray {
		return 0, &TypeMismatchError{Expected: StateStartArray, Actual: state}
	}

	if len(r.nestingStack) >= r.maxNestingDepth {
		return 0, ErrNestingDepthExceeded
	}

	r.invalidateState()

	if r.data[r.offset] == encodeInitialByte(MajorTypeArray, byte(AdditionalInfoIndefiniteLength)) {
		if r.conformanceMode >= ConformanceCanonical {
			return 0, ErrIndefiniteLengthNotAllowed
		}
		r.offset++
		r.nestingStack = append(r.nestingStack, readerNestingInfo{
			majorType:      MajorTypeArray,
			definiteLength: -1,
			isIndefinite:   true,
		})
		return -1, nil
	}

	length, err := r.readArgumentValue(MajorTypeArray)
	if err != nil {
		return 0, err
	}

	r.nestingStack = append(r.nestingStack, readerNestingInfo{
		majorType:      MajorTypeArray,
		definiteLength: int64(length),
	})

	return int(length), nil
}

// ReadEndArray reads the end of an array.
func (r *CborReader) ReadEndArray() error {
	state, err := r.PeekState()
	if err != nil {
		return err
	}
	if state != StateEndArray {
		return &TypeMismatchError{Expected: StateEndArray, Actual: state}
	}

	if len(r.nestingStack) == 0 {
		return ErrInvalidState
	}

	info := &r.nestingStack[len(r.nestingStack)-1]
	if info.majorType != MajorTypeArray {
		return ErrInvalidState
	}

	if info.isIndefinite {
		if r.data[r.offset] != breakByte {
			return ErrMissingBreak
		}
		r.offset++
	}

	r.nestingStack = r.nestingStack[:len(r.nestingStack)-1]
	r.invalidateState()
	r.advanceContainer()
	return nil
}

// ReadStartMap reads the start of a map and returns its length.
// Returns -1 for indefinite-length maps.
func (r *CborReader) ReadStartMap() (int, error) {
	state, err := r.PeekState()
	if err != nil {
		return 0, err
	}
	if state != StateStartMap {
		return 0, &TypeMismatchError{Expected: StateStartMap, Actual: state}
	}

	if len(r.nestingStack) >= r.maxNestingDepth {
		return 0, ErrNestingDepthExceeded
	}

	r.invalidateState()

	if r.data[r.offset] == encodeInitialByte(MajorTypeMap, byte(AdditionalInfoIndefiniteLength)) {
		if r.conformanceMode >= ConformanceCanonical {
			return 0, ErrIndefiniteLengthNotAllowed
		}
		r.offset++
		r.nestingStack = append(r.nestingStack, readerNestingInfo{
			majorType:      MajorTypeMap,
			definiteLength: -1,
			isMap:          true,
			isIndefinite:   true,
		})
		return -1, nil
	}

	length, err := r.readArgumentValue(MajorTypeMap)
	if err != nil {
		return 0, err
	}

	r.nestingStack = append(r.nestingStack, readerNestingInfo{
		majorType:      MajorTypeMap,
		definiteLength: int64(length),
		isMap:          true,
	})

	return int(length), nil
}

// ReadEndMap reads the end of a map.
func (r *CborReader) ReadEndMap() error {
	state, err := r.PeekState()
	if err != nil {
		return err
	}
	if state != StateEndMap {
		return &TypeMismatchError{Expected: StateEndMap, Actual: state}
	}

	if len(r.nestingStack) == 0 {
		return ErrInvalidState
	}

	info := &r.nestingStack[len(r.nestingStack)-1]
	if info.majorType != MajorTypeMap {
		return ErrInvalidState
	}

	if info.isIndefinite {
		if r.data[r.offset] != breakByte {
			return ErrMissingBreak
		}
		r.offset++
	}

	r.nestingStack = r.nestingStack[:len(r.nestingStack)-1]
	r.invalidateState()
	r.advanceContainer()
	return nil
}

// ReadTag reads a semantic tag.
func (r *CborReader) ReadTag() (CborTag, error) {
	state, err := r.PeekState()
	if err != nil {
		return 0, err
	}
	if state != StateTag {
		return 0, &TypeMismatchError{Expected: StateTag, Actual: state}
	}

	r.invalidateState()
	val, err := r.readArgumentValue(MajorTypeTag)
	if err != nil {
		return 0, err
	}

	// Don't advance container - the tagged value will do that
	return CborTag(val), nil
}

// ReadBoolean reads a boolean value.
func (r *CborReader) ReadBoolean() (bool, error) {
	state, err := r.PeekState()
	if err != nil {
		return false, err
	}
	if state != StateBoolean {
		return false, &TypeMismatchError{Expected: StateBoolean, Actual: state}
	}

	r.invalidateState()
	_, ai := decodeInitialByte(r.data[r.offset])
	r.offset++
	r.advanceContainer()

	return ai == byte(SimpleValueTrue), nil
}

// ReadNull reads a null value.
func (r *CborReader) ReadNull() error {
	state, err := r.PeekState()
	if err != nil {
		return err
	}
	if state != StateNull {
		return &TypeMismatchError{Expected: StateNull, Actual: state}
	}

	r.invalidateState()
	r.offset++
	r.advanceContainer()
	return nil
}

// ReadUndefined reads an undefined value.
func (r *CborReader) ReadUndefined() error {
	state, err := r.PeekState()
	if err != nil {
		return err
	}
	if state != StateUndefinedValue {
		return &TypeMismatchError{Expected: StateUndefinedValue, Actual: state}
	}

	r.invalidateState()
	r.offset++
	r.advanceContainer()
	return nil
}

// ReadSimpleValue reads a simple value.
func (r *CborReader) ReadSimpleValue() (SimpleValue, error) {
	state, err := r.PeekState()
	if err != nil {
		return 0, err
	}

	// Allow reading booleans, null, undefined as simple values too
	switch state {
	case StateSimpleValue, StateBoolean, StateNull, StateUndefinedValue:
		// ok
	default:
		return 0, &TypeMismatchError{Expected: StateSimpleValue, Actual: state}
	}

	r.invalidateState()
	_, ai := decodeInitialByte(r.data[r.offset])
	r.offset++

	var value SimpleValue
	if ai == 24 {
		if r.offset >= len(r.data) {
			return 0, ErrUnexpectedEndOfData
		}
		value = SimpleValue(r.data[r.offset])
		r.offset++

		// Canonical check: value must be >= 32
		if r.conformanceMode >= ConformanceStrict && value < 32 {
			return 0, ErrNonCanonical
		}
	} else {
		value = SimpleValue(ai)
	}

	r.advanceContainer()
	return value, nil
}

// ReadFloat16 reads a half-precision floating-point number.
func (r *CborReader) ReadFloat16() (float32, error) {
	state, err := r.PeekState()
	if err != nil {
		return 0, err
	}
	if state != StateHalfPrecisionFloat {
		return 0, &TypeMismatchError{Expected: StateHalfPrecisionFloat, Actual: state}
	}

	r.invalidateState()
	r.offset++ // Skip initial byte

	if r.offset+2 > len(r.data) {
		return 0, ErrUnexpectedEndOfData
	}

	bits := binary.BigEndian.Uint16(r.data[r.offset:])
	r.offset += 2
	r.advanceContainer()

	return float16BitsToFloat32(bits), nil
}

// ReadFloat32 reads a single-precision floating-point number.
func (r *CborReader) ReadFloat32() (float32, error) {
	state, err := r.PeekState()
	if err != nil {
		return 0, err
	}
	if state != StateSinglePrecisionFloat {
		return 0, &TypeMismatchError{Expected: StateSinglePrecisionFloat, Actual: state}
	}

	r.invalidateState()
	r.offset++ // Skip initial byte

	if r.offset+4 > len(r.data) {
		return 0, ErrUnexpectedEndOfData
	}

	bits := binary.BigEndian.Uint32(r.data[r.offset:])
	r.offset += 4
	r.advanceContainer()

	return math.Float32frombits(bits), nil
}

// ReadFloat64 reads a double-precision floating-point number.
func (r *CborReader) ReadFloat64() (float64, error) {
	state, err := r.PeekState()
	if err != nil {
		return 0, err
	}
	if state != StateDoublePrecisionFloat {
		return 0, &TypeMismatchError{Expected: StateDoublePrecisionFloat, Actual: state}
	}

	r.invalidateState()
	r.offset++ // Skip initial byte

	if r.offset+8 > len(r.data) {
		return 0, ErrUnexpectedEndOfData
	}

	bits := binary.BigEndian.Uint64(r.data[r.offset:])
	r.offset += 8
	r.advanceContainer()

	return math.Float64frombits(bits), nil
}

// ReadFloat reads any floating-point number and returns it as float64.
func (r *CborReader) ReadFloat() (float64, error) {
	state, err := r.PeekState()
	if err != nil {
		return 0, err
	}

	switch state {
	case StateHalfPrecisionFloat:
		f, err := r.ReadFloat16()
		return float64(f), err
	case StateSinglePrecisionFloat:
		f, err := r.ReadFloat32()
		return float64(f), err
	case StateDoublePrecisionFloat:
		return r.ReadFloat64()
	default:
		return 0, &TypeMismatchError{Expected: StateDoublePrecisionFloat, Actual: state}
	}
}

// ReadDateTimeString reads a date/time string (tag 0).
func (r *CborReader) ReadDateTimeString() (time.Time, error) {
	tag, err := r.ReadTag()
	if err != nil {
		return time.Time{}, err
	}
	if tag != TagDateTimeString {
		return time.Time{}, NewCborError(ErrInvalidCbor, r.offset, "expected datetime string tag")
	}

	str, err := r.ReadTextString()
	if err != nil {
		return time.Time{}, err
	}

	return time.Parse(time.RFC3339Nano, str)
}

// ReadUnixTime reads an epoch-based date/time (tag 1).
func (r *CborReader) ReadUnixTime() (time.Time, error) {
	tag, err := r.ReadTag()
	if err != nil {
		return time.Time{}, err
	}
	if tag != TagUnixTime {
		return time.Time{}, NewCborError(ErrInvalidCbor, r.offset, "expected unix time tag")
	}

	state, err := r.PeekState()
	if err != nil {
		return time.Time{}, err
	}

	switch state {
	case StateUnsignedInteger, StateNegativeInteger:
		secs, err := r.ReadInt64()
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(secs, 0), nil

	case StateHalfPrecisionFloat, StateSinglePrecisionFloat, StateDoublePrecisionFloat:
		f, err := r.ReadFloat()
		if err != nil {
			return time.Time{}, err
		}
		secs := int64(f)
		nsecs := int64((f - float64(secs)) * 1e9)
		return time.Unix(secs, nsecs), nil

	default:
		return time.Time{}, &TypeMismatchError{Expected: StateUnsignedInteger, Actual: state}
	}
}

// SkipValue skips the current value (including nested values for arrays/maps).
func (r *CborReader) SkipValue() error {
	state, err := r.PeekState()
	if err != nil {
		return err
	}

	switch state {
	case StateUnsignedInteger:
		_, err = r.ReadUint64()
		return err
	case StateNegativeInteger:
		_, err = r.ReadInt64()
		return err
	case StateByteString, StateStartIndefiniteLengthByteString:
		_, err = r.ReadByteString()
		return err
	case StateTextString, StateStartIndefiniteLengthTextString:
		_, err = r.ReadTextString()
		return err
	case StateStartArray:
		return r.skipArray()
	case StateStartMap:
		return r.skipMap()
	case StateTag:
		_, err = r.ReadTag()
		if err != nil {
			return err
		}
		return r.SkipValue()
	case StateBoolean:
		_, err = r.ReadBoolean()
		return err
	case StateNull:
		return r.ReadNull()
	case StateUndefinedValue:
		return r.ReadUndefined()
	case StateSimpleValue:
		_, err = r.ReadSimpleValue()
		return err
	case StateHalfPrecisionFloat:
		_, err = r.ReadFloat16()
		return err
	case StateSinglePrecisionFloat:
		_, err = r.ReadFloat32()
		return err
	case StateDoublePrecisionFloat:
		_, err = r.ReadFloat64()
		return err
	default:
		return ErrInvalidState
	}
}

// skipArray skips an array and all its contents.
func (r *CborReader) skipArray() error {
	length, err := r.ReadStartArray()
	if err != nil {
		return err
	}

	if length == -1 {
		// Indefinite length
		for {
			state, err := r.PeekState()
			if err != nil {
				return err
			}
			if state == StateEndArray {
				break
			}
			if err := r.SkipValue(); err != nil {
				return err
			}
		}
	} else {
		for i := 0; i < length; i++ {
			if err := r.SkipValue(); err != nil {
				return err
			}
		}
	}

	return r.ReadEndArray()
}

// skipMap skips a map and all its contents.
func (r *CborReader) skipMap() error {
	length, err := r.ReadStartMap()
	if err != nil {
		return err
	}

	if length == -1 {
		// Indefinite length
		for {
			state, err := r.PeekState()
			if err != nil {
				return err
			}
			if state == StateEndMap {
				break
			}
			// Skip key
			if err := r.SkipValue(); err != nil {
				return err
			}
			// Skip value
			if err := r.SkipValue(); err != nil {
				return err
			}
		}
	} else {
		for i := 0; i < length; i++ {
			// Skip key
			if err := r.SkipValue(); err != nil {
				return err
			}
			// Skip value
			if err := r.SkipValue(); err != nil {
				return err
			}
		}
	}

	return r.ReadEndMap()
}

// TryReadNull returns true if the next value is null and consumes it.
func (r *CborReader) TryReadNull() (bool, error) {
	state, err := r.PeekState()
	if err != nil {
		return false, err
	}
	if state == StateNull {
		return true, r.ReadNull()
	}
	return false, nil
}

// ReadEncodedValue reads a single complete CBOR value as raw bytes.
func (r *CborReader) ReadEncodedValue() ([]byte, error) {
	start := r.offset
	err := r.SkipValue()
	if err != nil {
		return nil, err
	}

	result := make([]byte, r.offset-start)
	copy(result, r.data[start:r.offset])
	return result, nil
}
