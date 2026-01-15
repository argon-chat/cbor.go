package cbor

import (
	"encoding/hex"
	"testing"
)

// RFC 8949 Appendix A test vectors
func TestRFC8949Appendix(t *testing.T) {
	tests := []struct {
		name     string
		hex      string
		testFunc func(t *testing.T, data []byte)
	}{
		{
			name: "0",
			hex:  "00",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadUint64()
				if err != nil {
					t.Fatalf("ReadUint64 failed: %v", err)
				}
				if val != 0 {
					t.Errorf("got %d, want 0", val)
				}
			},
		},
		{
			name: "1",
			hex:  "01",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadUint64()
				if err != nil {
					t.Fatalf("ReadUint64 failed: %v", err)
				}
				if val != 1 {
					t.Errorf("got %d, want 1", val)
				}
			},
		},
		{
			name: "10",
			hex:  "0a",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadUint64()
				if err != nil {
					t.Fatalf("ReadUint64 failed: %v", err)
				}
				if val != 10 {
					t.Errorf("got %d, want 10", val)
				}
			},
		},
		{
			name: "23",
			hex:  "17",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadUint64()
				if err != nil {
					t.Fatalf("ReadUint64 failed: %v", err)
				}
				if val != 23 {
					t.Errorf("got %d, want 23", val)
				}
			},
		},
		{
			name: "24",
			hex:  "1818",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadUint64()
				if err != nil {
					t.Fatalf("ReadUint64 failed: %v", err)
				}
				if val != 24 {
					t.Errorf("got %d, want 24", val)
				}
			},
		},
		{
			name: "25",
			hex:  "1819",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadUint64()
				if err != nil {
					t.Fatalf("ReadUint64 failed: %v", err)
				}
				if val != 25 {
					t.Errorf("got %d, want 25", val)
				}
			},
		},
		{
			name: "100",
			hex:  "1864",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadUint64()
				if err != nil {
					t.Fatalf("ReadUint64 failed: %v", err)
				}
				if val != 100 {
					t.Errorf("got %d, want 100", val)
				}
			},
		},
		{
			name: "1000",
			hex:  "1903e8",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadUint64()
				if err != nil {
					t.Fatalf("ReadUint64 failed: %v", err)
				}
				if val != 1000 {
					t.Errorf("got %d, want 1000", val)
				}
			},
		},
		{
			name: "1000000",
			hex:  "1a000f4240",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadUint64()
				if err != nil {
					t.Fatalf("ReadUint64 failed: %v", err)
				}
				if val != 1000000 {
					t.Errorf("got %d, want 1000000", val)
				}
			},
		},
		{
			name: "1000000000000",
			hex:  "1b000000e8d4a51000",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadUint64()
				if err != nil {
					t.Fatalf("ReadUint64 failed: %v", err)
				}
				if val != 1000000000000 {
					t.Errorf("got %d, want 1000000000000", val)
				}
			},
		},
		{
			name: "-1",
			hex:  "20",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadInt64()
				if err != nil {
					t.Fatalf("ReadInt64 failed: %v", err)
				}
				if val != -1 {
					t.Errorf("got %d, want -1", val)
				}
			},
		},
		{
			name: "-10",
			hex:  "29",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadInt64()
				if err != nil {
					t.Fatalf("ReadInt64 failed: %v", err)
				}
				if val != -10 {
					t.Errorf("got %d, want -10", val)
				}
			},
		},
		{
			name: "-100",
			hex:  "3863",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadInt64()
				if err != nil {
					t.Fatalf("ReadInt64 failed: %v", err)
				}
				if val != -100 {
					t.Errorf("got %d, want -100", val)
				}
			},
		},
		{
			name: "-1000",
			hex:  "3903e7",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadInt64()
				if err != nil {
					t.Fatalf("ReadInt64 failed: %v", err)
				}
				if val != -1000 {
					t.Errorf("got %d, want -1000", val)
				}
			},
		},
		{
			name: "empty_byte_string",
			hex:  "40",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadByteString()
				if err != nil {
					t.Fatalf("ReadByteString failed: %v", err)
				}
				if len(val) != 0 {
					t.Errorf("got len %d, want 0", len(val))
				}
			},
		},
		{
			name: "h'01020304'",
			hex:  "4401020304",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadByteString()
				if err != nil {
					t.Fatalf("ReadByteString failed: %v", err)
				}
				expected := []byte{1, 2, 3, 4}
				for i, b := range val {
					if b != expected[i] {
						t.Errorf("byte %d: got %d, want %d", i, b, expected[i])
					}
				}
			},
		},
		{
			name: "empty_text_string",
			hex:  "60",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadTextString()
				if err != nil {
					t.Fatalf("ReadTextString failed: %v", err)
				}
				if val != "" {
					t.Errorf("got %q, want empty string", val)
				}
			},
		},
		{
			name: "a",
			hex:  "6161",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadTextString()
				if err != nil {
					t.Fatalf("ReadTextString failed: %v", err)
				}
				if val != "a" {
					t.Errorf("got %q, want 'a'", val)
				}
			},
		},
		{
			name: "IETF",
			hex:  "6449455446",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadTextString()
				if err != nil {
					t.Fatalf("ReadTextString failed: %v", err)
				}
				if val != "IETF" {
					t.Errorf("got %q, want 'IETF'", val)
				}
			},
		},
		{
			name: "backslash_quote",
			hex:  "62225c",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadTextString()
				if err != nil {
					t.Fatalf("ReadTextString failed: %v", err)
				}
				if val != "\"\\" {
					t.Errorf("got %q, want '\"\\\\'", val)
				}
			},
		},
		{
			name: "unicode_u",
			hex:  "62c3bc",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadTextString()
				if err != nil {
					t.Fatalf("ReadTextString failed: %v", err)
				}
				if val != "\u00fc" {
					t.Errorf("got %q, want 'ü'", val)
				}
			},
		},
		{
			name: "empty_array",
			hex:  "80",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				length, err := r.ReadStartArray()
				if err != nil {
					t.Fatalf("ReadStartArray failed: %v", err)
				}
				if length != 0 {
					t.Errorf("got length %d, want 0", length)
				}
				if err := r.ReadEndArray(); err != nil {
					t.Fatalf("ReadEndArray failed: %v", err)
				}
			},
		},
		{
			name: "[1, 2, 3]",
			hex:  "83010203",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				length, err := r.ReadStartArray()
				if err != nil {
					t.Fatalf("ReadStartArray failed: %v", err)
				}
				if length != 3 {
					t.Errorf("got length %d, want 3", length)
				}
				for i := int64(1); i <= 3; i++ {
					val, err := r.ReadInt64()
					if err != nil {
						t.Fatalf("ReadInt64 failed: %v", err)
					}
					if val != i {
						t.Errorf("got %d, want %d", val, i)
					}
				}
				if err := r.ReadEndArray(); err != nil {
					t.Fatalf("ReadEndArray failed: %v", err)
				}
			},
		},
		{
			name: "[[1], [2, 3], [4, 5]]",
			hex:  "83810182020382040500",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				length, err := r.ReadStartArray()
				if err != nil {
					t.Fatalf("ReadStartArray failed: %v", err)
				}
				if length != 3 {
					t.Errorf("got length %d, want 3", length)
				}
				// [1]
				l1, _ := r.ReadStartArray()
				if l1 != 1 {
					t.Errorf("got length %d, want 1", l1)
				}
				v1, _ := r.ReadInt64()
				if v1 != 1 {
					t.Errorf("got %d, want 1", v1)
				}
				_ = r.ReadEndArray()
				// [2, 3]
				l2, _ := r.ReadStartArray()
				if l2 != 2 {
					t.Errorf("got length %d, want 2", l2)
				}
				v2, _ := r.ReadInt64()
				if v2 != 2 {
					t.Errorf("got %d, want 2", v2)
				}
				v3, _ := r.ReadInt64()
				if v3 != 3 {
					t.Errorf("got %d, want 3", v3)
				}
				_ = r.ReadEndArray()
				// [4, 5]
				l3, _ := r.ReadStartArray()
				if l3 != 2 {
					t.Errorf("got length %d, want 2", l3)
				}
				v4, _ := r.ReadInt64()
				if v4 != 4 {
					t.Errorf("got %d, want 4", v4)
				}
				v5, _ := r.ReadInt64()
				if v5 != 5 {
					t.Errorf("got %d, want 5", v5)
				}
				_ = r.ReadEndArray()
				_ = r.ReadEndArray()
			},
		},
		{
			name: "empty_map",
			hex:  "a0",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				length, err := r.ReadStartMap()
				if err != nil {
					t.Fatalf("ReadStartMap failed: %v", err)
				}
				if length != 0 {
					t.Errorf("got length %d, want 0", length)
				}
				if err := r.ReadEndMap(); err != nil {
					t.Fatalf("ReadEndMap failed: %v", err)
				}
			},
		},
		{
			name: "{1: 2, 3: 4}",
			hex:  "a201020304",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				length, err := r.ReadStartMap()
				if err != nil {
					t.Fatalf("ReadStartMap failed: %v", err)
				}
				if length != 2 {
					t.Errorf("got length %d, want 2", length)
				}
				k1, _ := r.ReadInt64()
				v1, _ := r.ReadInt64()
				if k1 != 1 || v1 != 2 {
					t.Errorf("got %d: %d, want 1: 2", k1, v1)
				}
				k2, _ := r.ReadInt64()
				v2, _ := r.ReadInt64()
				if k2 != 3 || v2 != 4 {
					t.Errorf("got %d: %d, want 3: 4", k2, v2)
				}
				_ = r.ReadEndMap()
			},
		},
		{
			name: "{'a': 1, 'b': [2, 3]}",
			hex:  "a26161016162820203",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				length, err := r.ReadStartMap()
				if err != nil {
					t.Fatalf("ReadStartMap failed: %v", err)
				}
				if length != 2 {
					t.Errorf("got length %d, want 2", length)
				}
				k1, _ := r.ReadTextString()
				v1, _ := r.ReadInt64()
				if k1 != "a" || v1 != 1 {
					t.Errorf("got %s: %d, want a: 1", k1, v1)
				}
				k2, _ := r.ReadTextString()
				if k2 != "b" {
					t.Errorf("got key %s, want b", k2)
				}
				arrLen, _ := r.ReadStartArray()
				if arrLen != 2 {
					t.Errorf("got array length %d, want 2", arrLen)
				}
				av1, _ := r.ReadInt64()
				av2, _ := r.ReadInt64()
				if av1 != 2 || av2 != 3 {
					t.Errorf("got [%d, %d], want [2, 3]", av1, av2)
				}
				_ = r.ReadEndArray()
				_ = r.ReadEndMap()
			},
		},
		{
			name: "false",
			hex:  "f4",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadBoolean()
				if err != nil {
					t.Fatalf("ReadBoolean failed: %v", err)
				}
				if val != false {
					t.Errorf("got %v, want false", val)
				}
			},
		},
		{
			name: "true",
			hex:  "f5",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadBoolean()
				if err != nil {
					t.Fatalf("ReadBoolean failed: %v", err)
				}
				if val != true {
					t.Errorf("got %v, want true", val)
				}
			},
		},
		{
			name: "null",
			hex:  "f6",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				err := r.ReadNull()
				if err != nil {
					t.Fatalf("ReadNull failed: %v", err)
				}
			},
		},
		{
			name: "undefined",
			hex:  "f7",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				err := r.ReadUndefined()
				if err != nil {
					t.Fatalf("ReadUndefined failed: %v", err)
				}
			},
		},
		{
			name: "simple(16)",
			hex:  "f0",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadSimpleValue()
				if err != nil {
					t.Fatalf("ReadSimpleValue failed: %v", err)
				}
				if val != 16 {
					t.Errorf("got %d, want 16", val)
				}
			},
		},
		{
			name: "simple(255)",
			hex:  "f8ff",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadSimpleValue()
				if err != nil {
					t.Fatalf("ReadSimpleValue failed: %v", err)
				}
				if val != 255 {
					t.Errorf("got %d, want 255", val)
				}
			},
		},
		{
			name: "0.0_half",
			hex:  "f90000",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadFloat16()
				if err != nil {
					t.Fatalf("ReadFloat16 failed: %v", err)
				}
				if val != 0.0 {
					t.Errorf("got %v, want 0.0", val)
				}
			},
		},
		{
			name: "1.0_half",
			hex:  "f93c00",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadFloat16()
				if err != nil {
					t.Fatalf("ReadFloat16 failed: %v", err)
				}
				if val != 1.0 {
					t.Errorf("got %v, want 1.0", val)
				}
			},
		},
		{
			name: "1.5_half",
			hex:  "f93e00",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadFloat16()
				if err != nil {
					t.Fatalf("ReadFloat16 failed: %v", err)
				}
				if val != 1.5 {
					t.Errorf("got %v, want 1.5", val)
				}
			},
		},
		{
			name: "100000.0_single",
			hex:  "fa47c35000",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadFloat32()
				if err != nil {
					t.Fatalf("ReadFloat32 failed: %v", err)
				}
				if val != 100000.0 {
					t.Errorf("got %v, want 100000.0", val)
				}
			},
		},
		{
			name: "1.1_double",
			hex:  "fb3ff199999999999a",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadFloat64()
				if err != nil {
					t.Fatalf("ReadFloat64 failed: %v", err)
				}
				if val != 1.1 {
					t.Errorf("got %v, want 1.1", val)
				}
			},
		},
		{
			name: "tag_0_datetime",
			hex:  "c074323031332d30332d32315432303a30343a30305a",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				tag, err := r.ReadTag()
				if err != nil {
					t.Fatalf("ReadTag failed: %v", err)
				}
				if tag != TagDateTimeString {
					t.Errorf("got tag %d, want %d", tag, TagDateTimeString)
				}
				str, err := r.ReadTextString()
				if err != nil {
					t.Fatalf("ReadTextString failed: %v", err)
				}
				if str != "2013-03-21T20:04:00Z" {
					t.Errorf("got %q, want '2013-03-21T20:04:00Z'", str)
				}
			},
		},
		{
			name: "tag_1_epoch",
			hex:  "c11a514b67b0",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				tag, err := r.ReadTag()
				if err != nil {
					t.Fatalf("ReadTag failed: %v", err)
				}
				if tag != TagUnixTime {
					t.Errorf("got tag %d, want %d", tag, TagUnixTime)
				}
				val, err := r.ReadUint64()
				if err != nil {
					t.Fatalf("ReadUint64 failed: %v", err)
				}
				if val != 1363896240 {
					t.Errorf("got %d, want 1363896240", val)
				}
			},
		},
		{
			name: "tag_32_uri",
			hex:  "d82076687474703a2f2f7777772e6578616d706c652e636f6d",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				tag, err := r.ReadTag()
				if err != nil {
					t.Fatalf("ReadTag failed: %v", err)
				}
				if tag != TagURI {
					t.Errorf("got tag %d, want %d", tag, TagURI)
				}
				str, err := r.ReadTextString()
				if err != nil {
					t.Fatalf("ReadTextString failed: %v", err)
				}
				if str != "http://www.example.com" {
					t.Errorf("got %q, want 'http://www.example.com'", str)
				}
			},
		},
		{
			name: "indefinite_byte_string",
			hex:  "5f42010243030405ff",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadByteString()
				if err != nil {
					t.Fatalf("ReadByteString failed: %v", err)
				}
				expected := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
				if len(val) != len(expected) {
					t.Errorf("got length %d, want %d", len(val), len(expected))
				}
				for i, b := range val {
					if b != expected[i] {
						t.Errorf("byte %d: got %d, want %d", i, b, expected[i])
					}
				}
			},
		},
		{
			name: "indefinite_text_string",
			hex:  "7f657374726561646d696e67ff",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				val, err := r.ReadTextString()
				if err != nil {
					t.Fatalf("ReadTextString failed: %v", err)
				}
				if val != "streaming" {
					t.Errorf("got %q, want 'streaming'", val)
				}
			},
		},
		{
			name: "indefinite_array",
			hex:  "9f018202039f0405ffff",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				length, err := r.ReadStartArray()
				if err != nil {
					t.Fatalf("ReadStartArray failed: %v", err)
				}
				if length != -1 {
					t.Errorf("got length %d, want -1 (indefinite)", length)
				}
				// Read 1
				v1, _ := r.ReadInt64()
				if v1 != 1 {
					t.Errorf("got %d, want 1", v1)
				}
				// Read [2, 3]
				arrLen, _ := r.ReadStartArray()
				if arrLen != 2 {
					t.Errorf("got array length %d, want 2", arrLen)
				}
				_, _ = r.ReadInt64()
				_, _ = r.ReadInt64()
				_ = r.ReadEndArray()
				// Read indefinite [4, 5]
				arrLen2, _ := r.ReadStartArray()
				if arrLen2 != -1 {
					t.Errorf("got array length %d, want -1", arrLen2)
				}
				_, _ = r.ReadInt64()
				_, _ = r.ReadInt64()
				_ = r.ReadEndArray()
				// End outer array
				_ = r.ReadEndArray()
			},
		},
		{
			name: "indefinite_map",
			hex:  "bf61610161629f0203ffff",
			testFunc: func(t *testing.T, data []byte) {
				r := NewCborReader(data)
				length, err := r.ReadStartMap()
				if err != nil {
					t.Fatalf("ReadStartMap failed: %v", err)
				}
				if length != -1 {
					t.Errorf("got length %d, want -1 (indefinite)", length)
				}
				// Read "a": 1
				k1, _ := r.ReadTextString()
				v1, _ := r.ReadInt64()
				if k1 != "a" || v1 != 1 {
					t.Errorf("got %s: %d, want a: 1", k1, v1)
				}
				// Read "b": [2, 3]
				k2, _ := r.ReadTextString()
				if k2 != "b" {
					t.Errorf("got key %s, want b", k2)
				}
				arrLen, _ := r.ReadStartArray()
				if arrLen != -1 {
					t.Errorf("got array length %d, want -1", arrLen)
				}
				_, _ = r.ReadInt64()
				_, _ = r.ReadInt64()
				_ = r.ReadEndArray()
				// End map
				_ = r.ReadEndMap()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := hex.DecodeString(tt.hex)
			if err != nil {
				t.Fatalf("failed to decode hex: %v", err)
			}
			tt.testFunc(t, data)
		})
	}
}

// Test that writer produces correct CBOR for known values
func TestWriterProducesCorrectCBOR(t *testing.T) {
	tests := []struct {
		name      string
		writeFunc func(w *CborWriter) error
		expected  string
	}{
		{
			name:      "0",
			writeFunc: func(w *CborWriter) error { return w.WriteUint64(0) },
			expected:  "00",
		},
		{
			name:      "1",
			writeFunc: func(w *CborWriter) error { return w.WriteUint64(1) },
			expected:  "01",
		},
		{
			name:      "23",
			writeFunc: func(w *CborWriter) error { return w.WriteUint64(23) },
			expected:  "17",
		},
		{
			name:      "24",
			writeFunc: func(w *CborWriter) error { return w.WriteUint64(24) },
			expected:  "1818",
		},
		{
			name:      "100",
			writeFunc: func(w *CborWriter) error { return w.WriteUint64(100) },
			expected:  "1864",
		},
		{
			name:      "1000",
			writeFunc: func(w *CborWriter) error { return w.WriteUint64(1000) },
			expected:  "1903e8",
		},
		{
			name:      "-1",
			writeFunc: func(w *CborWriter) error { return w.WriteInt64(-1) },
			expected:  "20",
		},
		{
			name:      "-10",
			writeFunc: func(w *CborWriter) error { return w.WriteInt64(-10) },
			expected:  "29",
		},
		{
			name:      "-100",
			writeFunc: func(w *CborWriter) error { return w.WriteInt64(-100) },
			expected:  "3863",
		},
		{
			name:      "empty_byte_string",
			writeFunc: func(w *CborWriter) error { return w.WriteByteString([]byte{}) },
			expected:  "40",
		},
		{
			name:      "empty_text_string",
			writeFunc: func(w *CborWriter) error { return w.WriteTextString("") },
			expected:  "60",
		},
		{
			name:      "text_a",
			writeFunc: func(w *CborWriter) error { return w.WriteTextString("a") },
			expected:  "6161",
		},
		{
			name: "empty_array",
			writeFunc: func(w *CborWriter) error {
				if err := w.WriteStartArray(0); err != nil {
					return err
				}
				return w.WriteEndArray()
			},
			expected: "80",
		},
		{
			name: "empty_map",
			writeFunc: func(w *CborWriter) error {
				if err := w.WriteStartMap(0); err != nil {
					return err
				}
				return w.WriteEndMap()
			},
			expected: "a0",
		},
		{
			name:      "false",
			writeFunc: func(w *CborWriter) error { return w.WriteBoolean(false) },
			expected:  "f4",
		},
		{
			name:      "true",
			writeFunc: func(w *CborWriter) error { return w.WriteBoolean(true) },
			expected:  "f5",
		},
		{
			name:      "null",
			writeFunc: func(w *CborWriter) error { return w.WriteNull() },
			expected:  "f6",
		},
		{
			name:      "undefined",
			writeFunc: func(w *CborWriter) error { return w.WriteUndefined() },
			expected:  "f7",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewCborWriter()
			if err := tt.writeFunc(w); err != nil {
				t.Fatalf("write failed: %v", err)
			}
			got := hex.EncodeToString(w.Bytes())
			if got != tt.expected {
				t.Errorf("got %s, want %s", got, tt.expected)
			}
		})
	}
}
