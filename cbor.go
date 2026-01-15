// Package cbor provides CBOR (Concise Binary Object Representation) encoding and decoding
// as defined in RFC 8949. This implementation is inspired by .NET's System.Formats.Cbor.
package cbor

// MajorType represents the CBOR major type (3-bit value in the initial byte).
type MajorType byte

const (
	// MajorTypeUnsignedInteger represents unsigned integer (major type 0).
	MajorTypeUnsignedInteger MajorType = 0
	// MajorTypeNegativeInteger represents negative integer (major type 1).
	MajorTypeNegativeInteger MajorType = 1
	// MajorTypeByteString represents byte string (major type 2).
	MajorTypeByteString MajorType = 2
	// MajorTypeTextString represents UTF-8 text string (major type 3).
	MajorTypeTextString MajorType = 3
	// MajorTypeArray represents array of data items (major type 4).
	MajorTypeArray MajorType = 4
	// MajorTypeMap represents map of pairs of data items (major type 5).
	MajorTypeMap MajorType = 5
	// MajorTypeTag represents tagged data item (major type 6).
	MajorTypeTag MajorType = 6
	// MajorTypeSimpleOrFloat represents simple values and floats (major type 7).
	MajorTypeSimpleOrFloat MajorType = 7
)

// String returns the string representation of the major type.
func (mt MajorType) String() string {
	switch mt {
	case MajorTypeUnsignedInteger:
		return "UnsignedInteger"
	case MajorTypeNegativeInteger:
		return "NegativeInteger"
	case MajorTypeByteString:
		return "ByteString"
	case MajorTypeTextString:
		return "TextString"
	case MajorTypeArray:
		return "Array"
	case MajorTypeMap:
		return "Map"
	case MajorTypeTag:
		return "Tag"
	case MajorTypeSimpleOrFloat:
		return "SimpleOrFloat"
	default:
		return "Unknown"
	}
}

// AdditionalInfo represents the additional information in the initial byte.
type AdditionalInfo byte

const (
	// AdditionalInfoDirect means the value is encoded directly in the additional info (0-23).
	AdditionalInfoDirect AdditionalInfo = 0
	// AdditionalInfo8Bit means the following byte contains the value.
	AdditionalInfo8Bit AdditionalInfo = 24
	// AdditionalInfo16Bit means the following 2 bytes contain the value.
	AdditionalInfo16Bit AdditionalInfo = 25
	// AdditionalInfo32Bit means the following 4 bytes contain the value.
	AdditionalInfo32Bit AdditionalInfo = 26
	// AdditionalInfo64Bit means the following 8 bytes contain the value.
	AdditionalInfo64Bit AdditionalInfo = 27
	// AdditionalInfoIndefiniteLength means indefinite length (used for strings, arrays, maps).
	AdditionalInfoIndefiniteLength AdditionalInfo = 31
)

// SimpleValue represents CBOR simple values.
type SimpleValue byte

const (
	// SimpleValueFalse represents the boolean value false.
	SimpleValueFalse SimpleValue = 20
	// SimpleValueTrue represents the boolean value true.
	SimpleValueTrue SimpleValue = 21
	// SimpleValueNull represents a null value.
	SimpleValueNull SimpleValue = 22
	// SimpleValueUndefined represents an undefined value.
	SimpleValueUndefined SimpleValue = 23
)

// CborTag represents well-known CBOR semantic tags.
type CborTag uint64

const (
	// TagDateTimeString is a standard date/time string (RFC 3339).
	TagDateTimeString CborTag = 0
	// TagUnixTime is an epoch-based date/time.
	TagUnixTime CborTag = 1
	// TagUnsignedBignum is a positive bignum.
	TagUnsignedBignum CborTag = 2
	// TagNegativeBignum is a negative bignum.
	TagNegativeBignum CborTag = 3
	// TagDecimalFraction is a decimal fraction.
	TagDecimalFraction CborTag = 4
	// TagBigFloat is a bigfloat.
	TagBigFloat CborTag = 5
	// TagExpectedBase64URL is expected conversion to base64url encoding.
	TagExpectedBase64URL CborTag = 21
	// TagExpectedBase64 is expected conversion to base64 encoding.
	TagExpectedBase64 CborTag = 22
	// TagExpectedBase16 is expected conversion to base16 encoding.
	TagExpectedBase16 CborTag = 23
	// TagEncodedCborData is encoded CBOR data item.
	TagEncodedCborData CborTag = 24
	// TagURI is a URI (RFC 3986).
	TagURI CborTag = 32
	// TagBase64URL is a base64url encoded text.
	TagBase64URL CborTag = 33
	// TagBase64 is a base64 encoded text.
	TagBase64 CborTag = 34
	// TagRegularExpression is a regular expression (PCRE/ECMA262).
	TagRegularExpression CborTag = 35
	// TagMIMEMessage is a MIME message (RFC 2045).
	TagMIMEMessage CborTag = 36
	// TagSelfDescribedCbor is a self-described CBOR.
	TagSelfDescribedCbor CborTag = 55799
)

// CborReaderState represents the current state of the CBOR reader.
type CborReaderState int

const (
	// StateUndefined means the reader state is undefined.
	StateUndefined CborReaderState = iota
	// StateUnsignedInteger means an unsigned integer is next.
	StateUnsignedInteger
	// StateNegativeInteger means a negative integer is next.
	StateNegativeInteger
	// StateByteString means a byte string is next.
	StateByteString
	// StateTextString means a text string is next.
	StateTextString
	// StateStartArray means the start of an array is next.
	StateStartArray
	// StateEndArray means the end of an array is next.
	StateEndArray
	// StateStartMap means the start of a map is next.
	StateStartMap
	// StateEndMap means the end of a map is next.
	StateEndMap
	// StateTag means a semantic tag is next.
	StateTag
	// StateSimpleValue means a simple value is next.
	StateSimpleValue
	// StateHalfPrecisionFloat means a half-precision float is next.
	StateHalfPrecisionFloat
	// StateSinglePrecisionFloat means a single-precision float is next.
	StateSinglePrecisionFloat
	// StateDoublePrecisionFloat means a double-precision float is next.
	StateDoublePrecisionFloat
	// StateNull means a null value is next.
	StateNull
	// StateBoolean means a boolean value is next.
	StateBoolean
	// StateUndefinedValue means an undefined value is next.
	StateUndefinedValue
	// StateStartIndefiniteLengthByteString means the start of an indefinite-length byte string.
	StateStartIndefiniteLengthByteString
	// StateEndIndefiniteLengthByteString means the end of an indefinite-length byte string.
	StateEndIndefiniteLengthByteString
	// StateStartIndefiniteLengthTextString means the start of an indefinite-length text string.
	StateStartIndefiniteLengthTextString
	// StateEndIndefiniteLengthTextString means the end of an indefinite-length text string.
	StateEndIndefiniteLengthTextString
	// StateFinished means all CBOR data has been read.
	StateFinished
)

// String returns the string representation of the reader state.
func (s CborReaderState) String() string {
	switch s {
	case StateUndefined:
		return "Undefined"
	case StateUnsignedInteger:
		return "UnsignedInteger"
	case StateNegativeInteger:
		return "NegativeInteger"
	case StateByteString:
		return "ByteString"
	case StateTextString:
		return "TextString"
	case StateStartArray:
		return "StartArray"
	case StateEndArray:
		return "EndArray"
	case StateStartMap:
		return "StartMap"
	case StateEndMap:
		return "EndMap"
	case StateTag:
		return "Tag"
	case StateSimpleValue:
		return "SimpleValue"
	case StateHalfPrecisionFloat:
		return "HalfPrecisionFloat"
	case StateSinglePrecisionFloat:
		return "SinglePrecisionFloat"
	case StateDoublePrecisionFloat:
		return "DoublePrecisionFloat"
	case StateNull:
		return "Null"
	case StateBoolean:
		return "Boolean"
	case StateUndefinedValue:
		return "Undefined"
	case StateStartIndefiniteLengthByteString:
		return "StartIndefiniteLengthByteString"
	case StateEndIndefiniteLengthByteString:
		return "EndIndefiniteLengthByteString"
	case StateStartIndefiniteLengthTextString:
		return "StartIndefiniteLengthTextString"
	case StateEndIndefiniteLengthTextString:
		return "EndIndefiniteLengthTextString"
	case StateFinished:
		return "Finished"
	default:
		return "Unknown"
	}
}

// CborConformanceMode specifies the conformance mode for CBOR operations.
type CborConformanceMode int

const (
	// ConformanceLax allows non-conforming CBOR data.
	ConformanceLax CborConformanceMode = iota
	// ConformanceStrict requires strict conformance to RFC 8949.
	ConformanceStrict
	// ConformanceCanonical requires canonical CBOR encoding (RFC 8949 Section 4.2.1).
	ConformanceCanonical
	// ConformanceCtap2Canonical requires CTAP2 canonical CBOR encoding.
	ConformanceCtap2Canonical
)

// Break byte used to terminate indefinite-length items.
const breakByte byte = 0xFF

// encodeInitialByte creates the initial byte from major type and additional info.
func encodeInitialByte(mt MajorType, ai byte) byte {
	return byte(mt)<<5 | (ai & 0x1F)
}

// decodeInitialByte extracts major type and additional info from initial byte.
func decodeInitialByte(b byte) (MajorType, byte) {
	return MajorType(b >> 5), b & 0x1F
}
