package cbor

import (
	"encoding/binary"
	"math"
	"math/big"
	"time"
)

// CborWriter provides methods for writing CBOR encoded data.
type CborWriter struct {
	buffer                  []byte
	conformanceMode         CborConformanceMode
	nestingStack            []nestingInfo
	maxNestingDepth         int
	currentOffset           int
	allowMultipleRootValues bool
	rootValueWritten        bool
}

// nestingInfo tracks the state of nested containers.
type nestingInfo struct {
	majorType      MajorType
	definiteLength int64 // -1 for indefinite
	itemsWritten   int64
	isMap          bool
	keyWritten     bool // for maps, tracks if we're expecting a value
	isIndefinite   bool
}

// WriterOption is a function that configures a CborWriter.
type WriterOption func(*CborWriter)

// WithConformanceMode sets the conformance mode for the writer.
func WithConformanceMode(mode CborConformanceMode) WriterOption {
	return func(w *CborWriter) {
		w.conformanceMode = mode
	}
}

// WithInitialCapacity sets the initial buffer capacity.
func WithInitialCapacity(capacity int) WriterOption {
	return func(w *CborWriter) {
		w.buffer = make([]byte, 0, capacity)
	}
}

// WithMaxNestingDepth sets the maximum nesting depth.
func WithMaxNestingDepth(depth int) WriterOption {
	return func(w *CborWriter) {
		w.maxNestingDepth = depth
	}
}

// WithAllowMultipleRootValues allows writing multiple root-level values.
func WithAllowMultipleRootValues(allow bool) WriterOption {
	return func(w *CborWriter) {
		w.allowMultipleRootValues = allow
	}
}

// NewCborWriter creates a new CborWriter with the specified options.
func NewCborWriter(opts ...WriterOption) *CborWriter {
	w := &CborWriter{
		buffer:          make([]byte, 0, 256),
		conformanceMode: ConformanceLax,
		nestingStack:    make([]nestingInfo, 0, 16),
		maxNestingDepth: 64,
	}

	for _, opt := range opts {
		opt(w)
	}

	return w
}

// Reset clears the writer for reuse.
func (w *CborWriter) Reset() {
	w.buffer = w.buffer[:0]
	w.nestingStack = w.nestingStack[:0]
	w.currentOffset = 0
	w.rootValueWritten = false
}

// Bytes returns the encoded CBOR data.
func (w *CborWriter) Bytes() []byte {
	return w.buffer
}

// BytesCopy returns a copy of the encoded CBOR data.
func (w *CborWriter) BytesCopy() []byte {
	result := make([]byte, len(w.buffer))
	copy(result, w.buffer)
	return result
}

// Len returns the current length of the encoded data.
func (w *CborWriter) Len() int {
	return len(w.buffer)
}

// NestingDepth returns the current nesting depth.
func (w *CborWriter) NestingDepth() int {
	return len(w.nestingStack)
}

// checkNestingDepth ensures we don't exceed the maximum nesting depth.
func (w *CborWriter) checkNestingDepth() error {
	if len(w.nestingStack) >= w.maxNestingDepth {
		return ErrNestingDepthExceeded
	}
	return nil
}

// advanceContainer updates container state after writing an item.
func (w *CborWriter) advanceContainer() {
	if len(w.nestingStack) == 0 {
		w.rootValueWritten = true
		return
	}

	info := &w.nestingStack[len(w.nestingStack)-1]
	if info.isMap {
		if info.keyWritten {
			// We just wrote a value
			info.keyWritten = false
			info.itemsWritten++
		} else {
			// We just wrote a key
			info.keyWritten = true
		}
	} else {
		info.itemsWritten++
	}
}

// writeInitialByte writes the initial byte for a data item.
func (w *CborWriter) writeInitialByte(mt MajorType, value uint64) {
	if value < 24 {
		w.buffer = append(w.buffer, encodeInitialByte(mt, byte(value)))
	} else if value <= math.MaxUint8 {
		w.buffer = append(w.buffer, encodeInitialByte(mt, byte(AdditionalInfo8Bit)), byte(value))
	} else if value <= math.MaxUint16 {
		w.buffer = append(w.buffer, encodeInitialByte(mt, byte(AdditionalInfo16Bit)))
		w.buffer = binary.BigEndian.AppendUint16(w.buffer, uint16(value))
	} else if value <= math.MaxUint32 {
		w.buffer = append(w.buffer, encodeInitialByte(mt, byte(AdditionalInfo32Bit)))
		w.buffer = binary.BigEndian.AppendUint32(w.buffer, uint32(value))
	} else {
		w.buffer = append(w.buffer, encodeInitialByte(mt, byte(AdditionalInfo64Bit)))
		w.buffer = binary.BigEndian.AppendUint64(w.buffer, value)
	}
	w.currentOffset = len(w.buffer)
}

// writeMinimalInitialByte writes the initial byte using minimal encoding (for canonical mode).
func (w *CborWriter) writeMinimalInitialByte(mt MajorType, value uint64) {
	if value < 24 {
		w.buffer = append(w.buffer, encodeInitialByte(mt, byte(value)))
	} else if value <= math.MaxUint8 {
		w.buffer = append(w.buffer, encodeInitialByte(mt, byte(AdditionalInfo8Bit)), byte(value))
	} else if value <= math.MaxUint16 {
		w.buffer = append(w.buffer, encodeInitialByte(mt, byte(AdditionalInfo16Bit)))
		w.buffer = binary.BigEndian.AppendUint16(w.buffer, uint16(value))
	} else if value <= math.MaxUint32 {
		w.buffer = append(w.buffer, encodeInitialByte(mt, byte(AdditionalInfo32Bit)))
		w.buffer = binary.BigEndian.AppendUint32(w.buffer, uint32(value))
	} else {
		w.buffer = append(w.buffer, encodeInitialByte(mt, byte(AdditionalInfo64Bit)))
		w.buffer = binary.BigEndian.AppendUint64(w.buffer, value)
	}
	w.currentOffset = len(w.buffer)
}

// WriteInt64 writes a signed 64-bit integer.
func (w *CborWriter) WriteInt64(value int64) error {
	if value >= 0 {
		w.writeMinimalInitialByte(MajorTypeUnsignedInteger, uint64(value))
	} else {
		// CBOR encodes negative integers as -1 - n, so the encoded value is -(value+1)
		w.writeMinimalInitialByte(MajorTypeNegativeInteger, uint64(-1-value))
	}
	w.advanceContainer()
	return nil
}

// WriteUint64 writes an unsigned 64-bit integer.
func (w *CborWriter) WriteUint64(value uint64) error {
	w.writeMinimalInitialByte(MajorTypeUnsignedInteger, value)
	w.advanceContainer()
	return nil
}

// WriteInt32 writes a signed 32-bit integer.
func (w *CborWriter) WriteInt32(value int32) error {
	return w.WriteInt64(int64(value))
}

// WriteUint32 writes an unsigned 32-bit integer.
func (w *CborWriter) WriteUint32(value uint32) error {
	return w.WriteUint64(uint64(value))
}

// WriteInt16 writes a signed 16-bit integer.
func (w *CborWriter) WriteInt16(value int16) error {
	return w.WriteInt64(int64(value))
}

// WriteUint16 writes an unsigned 16-bit integer.
func (w *CborWriter) WriteUint16(value uint16) error {
	return w.WriteUint64(uint64(value))
}

// WriteInt8 writes a signed 8-bit integer.
func (w *CborWriter) WriteInt8(value int8) error {
	return w.WriteInt64(int64(value))
}

// WriteUint8 writes an unsigned 8-bit integer.
func (w *CborWriter) WriteUint8(value uint8) error {
	return w.WriteUint64(uint64(value))
}

// WriteInt writes an int value.
func (w *CborWriter) WriteInt(value int) error {
	return w.WriteInt64(int64(value))
}

// WriteBigInt writes a big integer using semantic tags 2 or 3.
func (w *CborWriter) WriteBigInt(value *big.Int) error {
	if value == nil {
		return w.WriteNull()
	}

	// Check if it fits in int64/uint64
	if value.IsInt64() {
		return w.WriteInt64(value.Int64())
	}
	if value.IsUint64() {
		return w.WriteUint64(value.Uint64())
	}

	// Need to use bignum encoding
	var tag CborTag
	var absValue *big.Int

	if value.Sign() >= 0 {
		tag = TagUnsignedBignum
		absValue = value
	} else {
		tag = TagNegativeBignum
		// For negative bignums, encode -(n+1) = -n - 1
		absValue = new(big.Int).Neg(value)
		absValue.Sub(absValue, big.NewInt(1))
	}

	// Write the tag
	if err := w.WriteTag(tag); err != nil {
		return err
	}

	// Write as byte string
	bytes := absValue.Bytes()
	return w.WriteByteString(bytes)
}

// WriteByteString writes a byte string.
func (w *CborWriter) WriteByteString(value []byte) error {
	w.writeMinimalInitialByte(MajorTypeByteString, uint64(len(value)))
	w.buffer = append(w.buffer, value...)
	w.currentOffset = len(w.buffer)
	w.advanceContainer()
	return nil
}

// WriteTextString writes a UTF-8 text string.
func (w *CborWriter) WriteTextString(value string) error {
	w.writeMinimalInitialByte(MajorTypeTextString, uint64(len(value)))
	w.buffer = append(w.buffer, value...)
	w.currentOffset = len(w.buffer)
	w.advanceContainer()
	return nil
}

// WriteStartArray writes the beginning of a definite-length array.
func (w *CborWriter) WriteStartArray(length int) error {
	if err := w.checkNestingDepth(); err != nil {
		return err
	}

	w.writeMinimalInitialByte(MajorTypeArray, uint64(length))
	w.nestingStack = append(w.nestingStack, nestingInfo{
		majorType:      MajorTypeArray,
		definiteLength: int64(length),
		isMap:          false,
		isIndefinite:   false,
	})
	return nil
}

// WriteStartIndefiniteLengthArray writes the beginning of an indefinite-length array.
func (w *CborWriter) WriteStartIndefiniteLengthArray() error {
	if w.conformanceMode == ConformanceCanonical || w.conformanceMode == ConformanceCtap2Canonical {
		return ErrIndefiniteLengthNotAllowed
	}

	if err := w.checkNestingDepth(); err != nil {
		return err
	}

	w.buffer = append(w.buffer, encodeInitialByte(MajorTypeArray, byte(AdditionalInfoIndefiniteLength)))
	w.currentOffset = len(w.buffer)
	w.nestingStack = append(w.nestingStack, nestingInfo{
		majorType:      MajorTypeArray,
		definiteLength: -1,
		isMap:          false,
		isIndefinite:   true,
	})
	return nil
}

// WriteEndArray writes the end of an array.
func (w *CborWriter) WriteEndArray() error {
	if len(w.nestingStack) == 0 {
		return ErrInvalidState
	}

	info := &w.nestingStack[len(w.nestingStack)-1]
	if info.majorType != MajorTypeArray {
		return ErrInvalidState
	}

	if info.isIndefinite {
		w.buffer = append(w.buffer, breakByte)
		w.currentOffset = len(w.buffer)
	} else if info.itemsWritten != info.definiteLength {
		if info.itemsWritten < info.definiteLength {
			return ErrIncompleteContainer
		}
		return ErrExtraItems
	}

	w.nestingStack = w.nestingStack[:len(w.nestingStack)-1]
	w.advanceContainer()
	return nil
}

// WriteStartMap writes the beginning of a definite-length map.
func (w *CborWriter) WriteStartMap(length int) error {
	if err := w.checkNestingDepth(); err != nil {
		return err
	}

	w.writeMinimalInitialByte(MajorTypeMap, uint64(length))
	w.nestingStack = append(w.nestingStack, nestingInfo{
		majorType:      MajorTypeMap,
		definiteLength: int64(length),
		isMap:          true,
		isIndefinite:   false,
	})
	return nil
}

// WriteStartIndefiniteLengthMap writes the beginning of an indefinite-length map.
func (w *CborWriter) WriteStartIndefiniteLengthMap() error {
	if w.conformanceMode == ConformanceCanonical || w.conformanceMode == ConformanceCtap2Canonical {
		return ErrIndefiniteLengthNotAllowed
	}

	if err := w.checkNestingDepth(); err != nil {
		return err
	}

	w.buffer = append(w.buffer, encodeInitialByte(MajorTypeMap, byte(AdditionalInfoIndefiniteLength)))
	w.currentOffset = len(w.buffer)
	w.nestingStack = append(w.nestingStack, nestingInfo{
		majorType:      MajorTypeMap,
		definiteLength: -1,
		isMap:          true,
		isIndefinite:   true,
	})
	return nil
}

// WriteEndMap writes the end of a map.
func (w *CborWriter) WriteEndMap() error {
	if len(w.nestingStack) == 0 {
		return ErrInvalidState
	}

	info := &w.nestingStack[len(w.nestingStack)-1]
	if info.majorType != MajorTypeMap {
		return ErrInvalidState
	}

	if info.keyWritten {
		// We wrote a key but no value
		return ErrIncompleteContainer
	}

	if info.isIndefinite {
		w.buffer = append(w.buffer, breakByte)
		w.currentOffset = len(w.buffer)
	} else if info.itemsWritten != info.definiteLength {
		if info.itemsWritten < info.definiteLength {
			return ErrIncompleteContainer
		}
		return ErrExtraItems
	}

	w.nestingStack = w.nestingStack[:len(w.nestingStack)-1]
	w.advanceContainer()
	return nil
}

// WriteTag writes a semantic tag.
func (w *CborWriter) WriteTag(tag CborTag) error {
	w.writeMinimalInitialByte(MajorTypeTag, uint64(tag))
	// Don't advance container - the tagged value will do that
	return nil
}

// WriteBoolean writes a boolean value.
func (w *CborWriter) WriteBoolean(value bool) error {
	if value {
		w.buffer = append(w.buffer, encodeInitialByte(MajorTypeSimpleOrFloat, byte(SimpleValueTrue)))
	} else {
		w.buffer = append(w.buffer, encodeInitialByte(MajorTypeSimpleOrFloat, byte(SimpleValueFalse)))
	}
	w.currentOffset = len(w.buffer)
	w.advanceContainer()
	return nil
}

// WriteNull writes a null value.
func (w *CborWriter) WriteNull() error {
	w.buffer = append(w.buffer, encodeInitialByte(MajorTypeSimpleOrFloat, byte(SimpleValueNull)))
	w.currentOffset = len(w.buffer)
	w.advanceContainer()
	return nil
}

// WriteUndefined writes an undefined value.
func (w *CborWriter) WriteUndefined() error {
	w.buffer = append(w.buffer, encodeInitialByte(MajorTypeSimpleOrFloat, byte(SimpleValueUndefined)))
	w.currentOffset = len(w.buffer)
	w.advanceContainer()
	return nil
}

// WriteSimpleValue writes a simple value.
func (w *CborWriter) WriteSimpleValue(value SimpleValue) error {
	if value < 32 {
		w.buffer = append(w.buffer, encodeInitialByte(MajorTypeSimpleOrFloat, byte(value)))
	} else {
		w.buffer = append(w.buffer, encodeInitialByte(MajorTypeSimpleOrFloat, byte(AdditionalInfo8Bit)), byte(value))
	}
	w.currentOffset = len(w.buffer)
	w.advanceContainer()
	return nil
}

// WriteFloat16 writes a half-precision (16-bit) floating-point number.
func (w *CborWriter) WriteFloat16(value float32) error {
	bits := float32ToFloat16Bits(value)
	w.buffer = append(w.buffer, encodeInitialByte(MajorTypeSimpleOrFloat, 25)) // 25 = half precision
	w.buffer = binary.BigEndian.AppendUint16(w.buffer, bits)
	w.currentOffset = len(w.buffer)
	w.advanceContainer()
	return nil
}

// WriteFloat32 writes a single-precision (32-bit) floating-point number.
func (w *CborWriter) WriteFloat32(value float32) error {
	bits := math.Float32bits(value)
	w.buffer = append(w.buffer, encodeInitialByte(MajorTypeSimpleOrFloat, 26)) // 26 = single precision
	w.buffer = binary.BigEndian.AppendUint32(w.buffer, bits)
	w.currentOffset = len(w.buffer)
	w.advanceContainer()
	return nil
}

// WriteFloat64 writes a double-precision (64-bit) floating-point number.
func (w *CborWriter) WriteFloat64(value float64) error {
	bits := math.Float64bits(value)
	w.buffer = append(w.buffer, encodeInitialByte(MajorTypeSimpleOrFloat, 27)) // 27 = double precision
	w.buffer = binary.BigEndian.AppendUint64(w.buffer, bits)
	w.currentOffset = len(w.buffer)
	w.advanceContainer()
	return nil
}

// WriteFloat writes a floating-point number using the smallest representation that doesn't lose precision.
func (w *CborWriter) WriteFloat(value float64) error {
	// Check if it can be represented as float32 without loss
	f32 := float32(value)
	if float64(f32) == value {
		// Check if it can be represented as float16 without loss
		f16bits := float32ToFloat16Bits(f32)
		f16back := float16BitsToFloat32(f16bits)
		if f16back == f32 && !math.IsNaN(value) {
			return w.WriteFloat16(f32)
		}
		return w.WriteFloat32(f32)
	}
	return w.WriteFloat64(value)
}

// WriteStartIndefiniteLengthByteString writes the start of an indefinite-length byte string.
func (w *CborWriter) WriteStartIndefiniteLengthByteString() error {
	if w.conformanceMode == ConformanceCanonical || w.conformanceMode == ConformanceCtap2Canonical {
		return ErrIndefiniteLengthNotAllowed
	}

	if err := w.checkNestingDepth(); err != nil {
		return err
	}

	w.buffer = append(w.buffer, encodeInitialByte(MajorTypeByteString, byte(AdditionalInfoIndefiniteLength)))
	w.currentOffset = len(w.buffer)
	w.nestingStack = append(w.nestingStack, nestingInfo{
		majorType:      MajorTypeByteString,
		definiteLength: -1,
		isIndefinite:   true,
	})
	return nil
}

// WriteByteStringChunk writes a chunk of an indefinite-length byte string.
func (w *CborWriter) WriteByteStringChunk(value []byte) error {
	if len(w.nestingStack) == 0 {
		return ErrInvalidState
	}
	info := &w.nestingStack[len(w.nestingStack)-1]
	if info.majorType != MajorTypeByteString || !info.isIndefinite {
		return ErrInvalidState
	}

	w.writeMinimalInitialByte(MajorTypeByteString, uint64(len(value)))
	w.buffer = append(w.buffer, value...)
	w.currentOffset = len(w.buffer)
	return nil
}

// WriteEndIndefiniteLengthByteString writes the end of an indefinite-length byte string.
func (w *CborWriter) WriteEndIndefiniteLengthByteString() error {
	if len(w.nestingStack) == 0 {
		return ErrInvalidState
	}
	info := &w.nestingStack[len(w.nestingStack)-1]
	if info.majorType != MajorTypeByteString || !info.isIndefinite {
		return ErrInvalidState
	}

	w.buffer = append(w.buffer, breakByte)
	w.currentOffset = len(w.buffer)
	w.nestingStack = w.nestingStack[:len(w.nestingStack)-1]
	w.advanceContainer()
	return nil
}

// WriteStartIndefiniteLengthTextString writes the start of an indefinite-length text string.
func (w *CborWriter) WriteStartIndefiniteLengthTextString() error {
	if w.conformanceMode == ConformanceCanonical || w.conformanceMode == ConformanceCtap2Canonical {
		return ErrIndefiniteLengthNotAllowed
	}

	if err := w.checkNestingDepth(); err != nil {
		return err
	}

	w.buffer = append(w.buffer, encodeInitialByte(MajorTypeTextString, byte(AdditionalInfoIndefiniteLength)))
	w.currentOffset = len(w.buffer)
	w.nestingStack = append(w.nestingStack, nestingInfo{
		majorType:      MajorTypeTextString,
		definiteLength: -1,
		isIndefinite:   true,
	})
	return nil
}

// WriteTextStringChunk writes a chunk of an indefinite-length text string.
func (w *CborWriter) WriteTextStringChunk(value string) error {
	if len(w.nestingStack) == 0 {
		return ErrInvalidState
	}
	info := &w.nestingStack[len(w.nestingStack)-1]
	if info.majorType != MajorTypeTextString || !info.isIndefinite {
		return ErrInvalidState
	}

	w.writeMinimalInitialByte(MajorTypeTextString, uint64(len(value)))
	w.buffer = append(w.buffer, value...)
	w.currentOffset = len(w.buffer)
	return nil
}

// WriteEndIndefiniteLengthTextString writes the end of an indefinite-length text string.
func (w *CborWriter) WriteEndIndefiniteLengthTextString() error {
	if len(w.nestingStack) == 0 {
		return ErrInvalidState
	}
	info := &w.nestingStack[len(w.nestingStack)-1]
	if info.majorType != MajorTypeTextString || !info.isIndefinite {
		return ErrInvalidState
	}

	w.buffer = append(w.buffer, breakByte)
	w.currentOffset = len(w.buffer)
	w.nestingStack = w.nestingStack[:len(w.nestingStack)-1]
	w.advanceContainer()
	return nil
}

// WriteDateTimeString writes a date/time string with the appropriate tag.
func (w *CborWriter) WriteDateTimeString(t time.Time) error {
	if err := w.WriteTag(TagDateTimeString); err != nil {
		return err
	}
	return w.WriteTextString(t.Format(time.RFC3339Nano))
}

// WriteUnixTime writes an epoch-based date/time with the appropriate tag.
func (w *CborWriter) WriteUnixTime(t time.Time) error {
	if err := w.WriteTag(TagUnixTime); err != nil {
		return err
	}
	// Use float if we need sub-second precision
	if t.Nanosecond() != 0 {
		seconds := float64(t.Unix()) + float64(t.Nanosecond())/1e9
		return w.WriteFloat64(seconds)
	}
	return w.WriteInt64(t.Unix())
}

// WriteUri writes a URI with the appropriate tag.
func (w *CborWriter) WriteUri(uri string) error {
	if err := w.WriteTag(TagURI); err != nil {
		return err
	}
	return w.WriteTextString(uri)
}

// WriteEncodedCborData writes encoded CBOR data with the appropriate tag.
func (w *CborWriter) WriteEncodedCborData(data []byte) error {
	if err := w.WriteTag(TagEncodedCborData); err != nil {
		return err
	}
	return w.WriteByteString(data)
}

// WriteSelfDescribedCbor writes the self-described CBOR tag.
func (w *CborWriter) WriteSelfDescribedCbor() error {
	return w.WriteTag(TagSelfDescribedCbor)
}

// WriteRaw writes raw bytes directly to the buffer.
// Use with caution - this bypasses all encoding.
func (w *CborWriter) WriteRaw(data []byte) error {
	w.buffer = append(w.buffer, data...)
	w.currentOffset = len(w.buffer)
	return nil
}

// float32ToFloat16Bits converts a float32 to IEEE 754 half-precision bits.
func float32ToFloat16Bits(f float32) uint16 {
	bits := math.Float32bits(f)
	sign := uint16((bits >> 16) & 0x8000)
	exp := int((bits >> 23) & 0xFF)
	frac := bits & 0x7FFFFF

	switch {
	case exp == 0:
		// Zero or subnormal
		return sign
	case exp == 255:
		// Inf or NaN
		if frac == 0 {
			return sign | 0x7C00
		}
		return sign | 0x7C00 | uint16(frac>>13)
	case exp > 142:
		// Overflow to infinity
		return sign | 0x7C00
	case exp < 113:
		// Underflow to zero
		return sign
	default:
		// Normal number
		exp16 := exp - 127 + 15
		frac16 := frac >> 13
		return sign | uint16(exp16<<10) | uint16(frac16)
	}
}

// float16BitsToFloat32 converts IEEE 754 half-precision bits to float32.
func float16BitsToFloat32(bits uint16) float32 {
	sign := uint32(bits&0x8000) << 16
	exp := int(bits>>10) & 0x1F
	frac := uint32(bits & 0x3FF)

	switch {
	case exp == 0:
		if frac == 0 {
			return math.Float32frombits(sign)
		}
		// Subnormal
		for frac&0x400 == 0 {
			frac <<= 1
			exp--
		}
		exp++
		frac &= 0x3FF
		fallthrough
	case exp < 31:
		exp32 := uint32(exp - 15 + 127)
		return math.Float32frombits(sign | (exp32 << 23) | (frac << 13))
	default:
		// Inf or NaN
		if frac == 0 {
			return math.Float32frombits(sign | 0x7F800000)
		}
		return math.Float32frombits(sign | 0x7F800000 | (frac << 13))
	}
}
