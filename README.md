# CBOR.go

[![CI](https://github.com/argon-chat/cbor.go/actions/workflows/ci.yml/badge.svg)](https://github.com/argon-chat/cbor.go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/argon-chat/cbor.go.svg)](https://pkg.go.dev/github.com/argon-chat/cbor.go)
[![Go Report Card](https://goreportcard.com/badge/github.com/argon-chat/cbor.go)](https://goreportcard.com/report/github.com/argon-chat/cbor.go)
[![codecov](https://codecov.io/gh/argon-chat/cbor.go/branch/main/graph/badge.svg)](https://codecov.io/gh/argon-chat/cbor.go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A complete CBOR (Concise Binary Object Representation) implementation for Go, inspired by .NET's `System.Formats.Cbor`.

## Features

- **Full RFC 8949 compliance** - Supports all major CBOR types
- **Low-level API** - Fine-grained control over encoding/decoding
- **Multiple conformance modes** - Lax, Strict, Canonical, and CTAP2 Canonical
- **Indefinite-length support** - Arrays, maps, byte strings, and text strings
- **Semantic tags** - Date/time, bignum, URI, and more
- **Half-precision floats** - IEEE 754 half-precision (16-bit) floating-point
- **Big integers** - Arbitrary precision integers via `math/big`
- **Skip functionality** - Efficiently skip complex nested structures
- **Zero allocations** - Minimal memory allocations during encoding

## Installation

```bash
go get github.com/argon-chat/cbor.go
```

## Quick Start

### Writing CBOR

```go
package main

import (
    "fmt"
    cbor "github.com/argon-chat/cbor.go"
)

func main() {
    w := cbor.NewCborWriter()

    // Write a map with two entries
    w.WriteStartMap(2)
    
    w.WriteTextString("name")
    w.WriteTextString("Alice")
    
    w.WriteTextString("age")
    w.WriteInt64(30)
    
    w.WriteEndMap()

    // Get the encoded bytes
    data := w.Bytes()
    fmt.Printf("CBOR: %x\n", data)
}
```

### Reading CBOR

```go
package main

import (
    "fmt"
    cbor "github.com/argon-chat/cbor.go"
)

func main() {
    data := []byte{0xa2, 0x64, 0x6e, 0x61, 0x6d, 0x65, 0x65, 0x41, 
                   0x6c, 0x69, 0x63, 0x65, 0x63, 0x61, 0x67, 0x65, 0x18, 0x1e}

    r := cbor.NewCborReader(data)

    length, _ := r.ReadStartMap()
    fmt.Printf("Map with %d entries:\n", length)

    for i := 0; i < length; i++ {
        key, _ := r.ReadTextString()
        
        state, _ := r.PeekState()
        switch state {
        case cbor.StateTextString:
            value, _ := r.ReadTextString()
            fmt.Printf("  %s: %s\n", key, value)
        case cbor.StateUnsignedInteger:
            value, _ := r.ReadInt64()
            fmt.Printf("  %s: %d\n", key, value)
        }
    }
    r.ReadEndMap()
}
```

## Supported Types

### Major Types

| Type | Description | Writer Methods | Reader Methods |
|------|-------------|----------------|----------------|
| 0 | Unsigned Integer | `WriteUint64`, `WriteUint32`, etc. | `ReadUint64`, `ReadUint32`, etc. |
| 1 | Negative Integer | `WriteInt64`, `WriteInt32`, etc. | `ReadInt64`, `ReadInt32`, etc. |
| 2 | Byte String | `WriteByteString` | `ReadByteString` |
| 3 | Text String | `WriteTextString` | `ReadTextString` |
| 4 | Array | `WriteStartArray`, `WriteEndArray` | `ReadStartArray`, `ReadEndArray` |
| 5 | Map | `WriteStartMap`, `WriteEndMap` | `ReadStartMap`, `ReadEndMap` |
| 6 | Tag | `WriteTag` | `ReadTag` |
| 7 | Simple/Float | `WriteBoolean`, `WriteNull`, `WriteFloat64`, etc. | `ReadBoolean`, `ReadNull`, `ReadFloat64`, etc. |

### Simple Values

- `false` (20) - `WriteBoolean(false)` / `ReadBoolean()`
- `true` (21) - `WriteBoolean(true)` / `ReadBoolean()`
- `null` (22) - `WriteNull()` / `ReadNull()`
- `undefined` (23) - `WriteUndefined()` / `ReadUndefined()`

### Floating-Point

- Half-precision (16-bit): `WriteFloat16` / `ReadFloat16`
- Single-precision (32-bit): `WriteFloat32` / `ReadFloat32`
- Double-precision (64-bit): `WriteFloat64` / `ReadFloat64`
- Auto-select smallest: `WriteFloat`

### Semantic Tags

| Tag | Description | Writer Method | Reader Method |
|-----|-------------|---------------|---------------|
| 0 | DateTime String (RFC 3339) | `WriteDateTimeString` | `ReadDateTimeString` |
| 1 | Unix Epoch Time | `WriteUnixTime` | `ReadUnixTime` |
| 2 | Positive Bignum | `WriteBigInt` | `ReadBigInt` |
| 3 | Negative Bignum | `WriteBigInt` | `ReadBigInt` |
| 32 | URI | `WriteUri` | via `ReadTag` + `ReadTextString` |
| 55799 | Self-Described CBOR | `WriteSelfDescribedCbor` | via `ReadTag` |

## Conformance Modes

```go
// Lax mode (default) - accepts all valid CBOR
w := cbor.NewCborWriter()

// Strict mode - validates encoding rules
w := cbor.NewCborWriter(cbor.WithConformanceMode(cbor.ConformanceStrict))

// Canonical mode - RFC 8949 Section 4.2.1
w := cbor.NewCborWriter(cbor.WithConformanceMode(cbor.ConformanceCanonical))

// CTAP2 Canonical mode
w := cbor.NewCborWriter(cbor.WithConformanceMode(cbor.ConformanceCtap2Canonical))
```

## Indefinite-Length Items

```go
// Indefinite-length array
w.WriteStartIndefiniteLengthArray()
w.WriteInt64(1)
w.WriteInt64(2)
w.WriteInt64(3)
w.WriteEndArray()

// Indefinite-length byte string
w.WriteStartIndefiniteLengthByteString()
w.WriteByteStringChunk([]byte{1, 2, 3})
w.WriteByteStringChunk([]byte{4, 5})
w.WriteEndIndefiniteLengthByteString()
```

## Advanced Usage

### Skipping Values

```go
r := cbor.NewCborReader(data)
length, _ := r.ReadStartArray()

// Skip the first element (can be complex nested structure)
r.SkipValue()

// Read the second element
value, _ := r.ReadInt64()
```

### Peeking State

```go
r := cbor.NewCborReader(data)

state, _ := r.PeekState()
switch state {
case cbor.StateTextString:
    str, _ := r.ReadTextString()
case cbor.StateUnsignedInteger:
    num, _ := r.ReadUint64()
case cbor.StateNull:
    r.ReadNull()
}
```

### Big Integers

```go
// Write a big integer
bigNum := new(big.Int)
bigNum.SetString("123456789012345678901234567890", 10)
w.WriteBigInt(bigNum)

// Read a big integer
r := cbor.NewCborReader(data)
bigNum, _ := r.ReadBigInt()
```

## Configuration Options

### Writer Options

- `WithConformanceMode(mode)` - Set conformance mode
- `WithInitialCapacity(size)` - Pre-allocate buffer
- `WithMaxNestingDepth(depth)` - Limit nesting depth (default: 64)
- `WithAllowMultipleRootValues(allow)` - Allow multiple root values

### Reader Options

- `WithReaderConformanceMode(mode)` - Set conformance mode
- `WithReaderMaxNestingDepth(depth)` - Limit nesting depth (default: 64)
- `WithReaderAllowMultipleRootValues(allow)` - Allow multiple root values

## Error Handling

The library provides detailed errors:

```go
val, err := r.ReadInt64()
if err != nil {
    if errors.Is(err, cbor.ErrUnexpectedEndOfData) {
        // Handle truncated data
    }
    if errors.Is(err, cbor.ErrOverflow) {
        // Handle integer overflow
    }
    if tmErr, ok := err.(*cbor.TypeMismatchError); ok {
        fmt.Printf("Expected %s but got %s\n", tmErr.Expected, tmErr.Actual)
    }
}
```

## Versioning

This project follows [Semantic Versioning](https://semver.org/):

- **MAJOR** version for incompatible API changes
- **MINOR** version for backwards-compatible functionality additions
- **PATCH** version for backwards-compatible bug fixes

You can check the version programmatically:

```go
fmt.Println(cbor.Version)      // "1.0.0"
fmt.Println(cbor.VersionInfo()) // "cbor.go v1.0.0"
```

### Creating a Release

```bash
# Tag a new release
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# The GitHub Action will automatically:
# - Run tests
# - Create a GitHub Release
# - Update the Go module proxy
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please ensure all tests pass and add tests for new functionality.

## License

MIT License
