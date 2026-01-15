# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2026-01-15

### Added

- Initial release
- `CborWriter` for encoding CBOR data
  - Unsigned and signed integers (int8 to int64, uint8 to uint64)
  - Big integers via `math/big.Int`
  - Byte strings and text strings
  - Definite-length arrays and maps
  - Indefinite-length arrays, maps, byte strings, and text strings
  - Semantic tags (datetime, unix time, bignum, URI, etc.)
  - IEEE 754 floating-point (half, single, double precision)
  - Simple values (true, false, null, undefined)
  - Self-described CBOR marker

- `CborReader` for decoding CBOR data
  - All types supported by CborWriter
  - State peeking without consuming data
  - Skip functionality for complex nested structures
  - Raw encoded value extraction

- Conformance modes
  - `ConformanceLax` - accepts all valid CBOR
  - `ConformanceStrict` - validates canonical encoding
  - `ConformanceCanonical` - RFC 8949 Section 4.2.1
  - `ConformanceCtap2Canonical` - CTAP2/FIDO2 compatibility

- Comprehensive error types
  - `CborError` with offset information
  - `TypeMismatchError` for type validation
  - Standard errors for common conditions

- Full RFC 8949 compliance with test vectors from Appendix A

[Unreleased]: https://github.com/argon-chat/cbor.go/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/argon-chat/cbor.go/releases/tag/v1.0.0
