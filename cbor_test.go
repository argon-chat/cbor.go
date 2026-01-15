package cbor

import (
	"bytes"
	"math"
	"math/big"
	"testing"
	"time"
)

func TestWriteReadUnsignedIntegers(t *testing.T) {
	tests := []struct {
		name  string
		value uint64
	}{
		{"zero", 0},
		{"one", 1},
		{"23", 23},
		{"24", 24},
		{"255", 255},
		{"256", 256},
		{"65535", 65535},
		{"65536", 65536},
		{"max_uint32", math.MaxUint32},
		{"max_uint32_plus_1", math.MaxUint32 + 1},
		{"max_uint64", math.MaxUint64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewCborWriter()
			if err := w.WriteUint64(tt.value); err != nil {
				t.Fatalf("WriteUint64 failed: %v", err)
			}

			r := NewCborReader(w.Bytes())
			got, err := r.ReadUint64()
			if err != nil {
				t.Fatalf("ReadUint64 failed: %v", err)
			}
			if got != tt.value {
				t.Errorf("got %d, want %d", got, tt.value)
			}
		})
	}
}

func TestWriteReadSignedIntegers(t *testing.T) {
	tests := []struct {
		name  string
		value int64
	}{
		{"zero", 0},
		{"one", 1},
		{"negative_one", -1},
		{"negative_24", -24},
		{"negative_25", -25},
		{"negative_256", -256},
		{"negative_257", -257},
		{"max_int64", math.MaxInt64},
		{"min_int64", math.MinInt64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewCborWriter()
			if err := w.WriteInt64(tt.value); err != nil {
				t.Fatalf("WriteInt64 failed: %v", err)
			}

			r := NewCborReader(w.Bytes())
			got, err := r.ReadInt64()
			if err != nil {
				t.Fatalf("ReadInt64 failed: %v", err)
			}
			if got != tt.value {
				t.Errorf("got %d, want %d", got, tt.value)
			}
		})
	}
}

func TestWriteReadByteString(t *testing.T) {
	tests := []struct {
		name  string
		value []byte
	}{
		{"empty", []byte{}},
		{"single_byte", []byte{0x01}},
		{"hello", []byte("hello")},
		{"long_string", bytes.Repeat([]byte{0xAB}, 1000)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewCborWriter()
			if err := w.WriteByteString(tt.value); err != nil {
				t.Fatalf("WriteByteString failed: %v", err)
			}

			r := NewCborReader(w.Bytes())
			got, err := r.ReadByteString()
			if err != nil {
				t.Fatalf("ReadByteString failed: %v", err)
			}
			if !bytes.Equal(got, tt.value) {
				t.Errorf("got %v, want %v", got, tt.value)
			}
		})
	}
}

func TestWriteReadTextString(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"empty", ""},
		{"hello", "hello"},
		{"unicode", "Привет мир! 🌍"},
		{"long_string", string(bytes.Repeat([]byte("a"), 1000))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewCborWriter()
			if err := w.WriteTextString(tt.value); err != nil {
				t.Fatalf("WriteTextString failed: %v", err)
			}

			r := NewCborReader(w.Bytes())
			got, err := r.ReadTextString()
			if err != nil {
				t.Fatalf("ReadTextString failed: %v", err)
			}
			if got != tt.value {
				t.Errorf("got %q, want %q", got, tt.value)
			}
		})
	}
}

func TestWriteReadBoolean(t *testing.T) {
	tests := []struct {
		name  string
		value bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewCborWriter()
			if err := w.WriteBoolean(tt.value); err != nil {
				t.Fatalf("WriteBoolean failed: %v", err)
			}

			r := NewCborReader(w.Bytes())
			got, err := r.ReadBoolean()
			if err != nil {
				t.Fatalf("ReadBoolean failed: %v", err)
			}
			if got != tt.value {
				t.Errorf("got %v, want %v", got, tt.value)
			}
		})
	}
}

func TestWriteReadNull(t *testing.T) {
	w := NewCborWriter()
	if err := w.WriteNull(); err != nil {
		t.Fatalf("WriteNull failed: %v", err)
	}

	r := NewCborReader(w.Bytes())
	if err := r.ReadNull(); err != nil {
		t.Fatalf("ReadNull failed: %v", err)
	}
}

func TestWriteReadUndefined(t *testing.T) {
	w := NewCborWriter()
	if err := w.WriteUndefined(); err != nil {
		t.Fatalf("WriteUndefined failed: %v", err)
	}

	r := NewCborReader(w.Bytes())
	if err := r.ReadUndefined(); err != nil {
		t.Fatalf("ReadUndefined failed: %v", err)
	}
}

func TestWriteReadFloat64(t *testing.T) {
	tests := []struct {
		name  string
		value float64
	}{
		{"zero", 0.0},
		{"one", 1.0},
		{"negative", -1.0},
		{"pi", 3.141592653589793},
		{"large", 1e100},
		{"small", 1e-100},
		{"inf", math.Inf(1)},
		{"neg_inf", math.Inf(-1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewCborWriter()
			if err := w.WriteFloat64(tt.value); err != nil {
				t.Fatalf("WriteFloat64 failed: %v", err)
			}

			r := NewCborReader(w.Bytes())
			got, err := r.ReadFloat64()
			if err != nil {
				t.Fatalf("ReadFloat64 failed: %v", err)
			}
			if got != tt.value {
				t.Errorf("got %v, want %v", got, tt.value)
			}
		})
	}
}

func TestWriteReadFloat64NaN(t *testing.T) {
	w := NewCborWriter()
	if err := w.WriteFloat64(math.NaN()); err != nil {
		t.Fatalf("WriteFloat64 failed: %v", err)
	}

	r := NewCborReader(w.Bytes())
	got, err := r.ReadFloat64()
	if err != nil {
		t.Fatalf("ReadFloat64 failed: %v", err)
	}
	if !math.IsNaN(got) {
		t.Errorf("got %v, want NaN", got)
	}
}

func TestWriteReadFloat32(t *testing.T) {
	tests := []struct {
		name  string
		value float32
	}{
		{"zero", 0.0},
		{"one", 1.0},
		{"negative", -1.0},
		{"pi", 3.1415927},
		{"inf", float32(math.Inf(1))},
		{"neg_inf", float32(math.Inf(-1))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewCborWriter()
			if err := w.WriteFloat32(tt.value); err != nil {
				t.Fatalf("WriteFloat32 failed: %v", err)
			}

			r := NewCborReader(w.Bytes())
			got, err := r.ReadFloat32()
			if err != nil {
				t.Fatalf("ReadFloat32 failed: %v", err)
			}
			if got != tt.value {
				t.Errorf("got %v, want %v", got, tt.value)
			}
		})
	}
}

func TestWriteReadFloat16(t *testing.T) {
	tests := []struct {
		name  string
		value float32
	}{
		{"zero", 0.0},
		{"one", 1.0},
		{"half", 0.5},
		{"inf", float32(math.Inf(1))},
		{"neg_inf", float32(math.Inf(-1))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewCborWriter()
			if err := w.WriteFloat16(tt.value); err != nil {
				t.Fatalf("WriteFloat16 failed: %v", err)
			}

			r := NewCborReader(w.Bytes())
			got, err := r.ReadFloat16()
			if err != nil {
				t.Fatalf("ReadFloat16 failed: %v", err)
			}
			if got != tt.value {
				t.Errorf("got %v, want %v", got, tt.value)
			}
		})
	}
}

func TestWriteReadArray(t *testing.T) {
	t.Run("empty_array", func(t *testing.T) {
		w := NewCborWriter()
		if err := w.WriteStartArray(0); err != nil {
			t.Fatalf("WriteStartArray failed: %v", err)
		}
		if err := w.WriteEndArray(); err != nil {
			t.Fatalf("WriteEndArray failed: %v", err)
		}

		r := NewCborReader(w.Bytes())
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
	})

	t.Run("array_with_integers", func(t *testing.T) {
		w := NewCborWriter()
		if err := w.WriteStartArray(3); err != nil {
			t.Fatalf("WriteStartArray failed: %v", err)
		}
		for _, v := range []int64{1, 2, 3} {
			if err := w.WriteInt64(v); err != nil {
				t.Fatalf("WriteInt64 failed: %v", err)
			}
		}
		if err := w.WriteEndArray(); err != nil {
			t.Fatalf("WriteEndArray failed: %v", err)
		}

		r := NewCborReader(w.Bytes())
		length, err := r.ReadStartArray()
		if err != nil {
			t.Fatalf("ReadStartArray failed: %v", err)
		}
		if length != 3 {
			t.Errorf("got length %d, want 3", length)
		}
		for _, expected := range []int64{1, 2, 3} {
			got, err := r.ReadInt64()
			if err != nil {
				t.Fatalf("ReadInt64 failed: %v", err)
			}
			if got != expected {
				t.Errorf("got %d, want %d", got, expected)
			}
		}
		if err := r.ReadEndArray(); err != nil {
			t.Fatalf("ReadEndArray failed: %v", err)
		}
	})

	t.Run("nested_arrays", func(t *testing.T) {
		w := NewCborWriter()
		if err := w.WriteStartArray(2); err != nil {
			t.Fatalf("WriteStartArray failed: %v", err)
		}
		if err := w.WriteStartArray(1); err != nil {
			t.Fatalf("WriteStartArray failed: %v", err)
		}
		if err := w.WriteInt64(1); err != nil {
			t.Fatalf("WriteInt64 failed: %v", err)
		}
		if err := w.WriteEndArray(); err != nil {
			t.Fatalf("WriteEndArray failed: %v", err)
		}
		if err := w.WriteStartArray(1); err != nil {
			t.Fatalf("WriteStartArray failed: %v", err)
		}
		if err := w.WriteInt64(2); err != nil {
			t.Fatalf("WriteInt64 failed: %v", err)
		}
		if err := w.WriteEndArray(); err != nil {
			t.Fatalf("WriteEndArray failed: %v", err)
		}
		if err := w.WriteEndArray(); err != nil {
			t.Fatalf("WriteEndArray failed: %v", err)
		}

		r := NewCborReader(w.Bytes())
		length, err := r.ReadStartArray()
		if err != nil {
			t.Fatalf("ReadStartArray failed: %v", err)
		}
		if length != 2 {
			t.Errorf("got length %d, want 2", length)
		}

		// First nested array
		innerLen, err := r.ReadStartArray()
		if err != nil {
			t.Fatalf("ReadStartArray failed: %v", err)
		}
		if innerLen != 1 {
			t.Errorf("got inner length %d, want 1", innerLen)
		}
		val, err := r.ReadInt64()
		if err != nil {
			t.Fatalf("ReadInt64 failed: %v", err)
		}
		if val != 1 {
			t.Errorf("got %d, want 1", val)
		}
		if err := r.ReadEndArray(); err != nil {
			t.Fatalf("ReadEndArray failed: %v", err)
		}

		// Second nested array
		innerLen, err = r.ReadStartArray()
		if err != nil {
			t.Fatalf("ReadStartArray failed: %v", err)
		}
		if innerLen != 1 {
			t.Errorf("got inner length %d, want 1", innerLen)
		}
		val, err = r.ReadInt64()
		if err != nil {
			t.Fatalf("ReadInt64 failed: %v", err)
		}
		if val != 2 {
			t.Errorf("got %d, want 2", val)
		}
		if err := r.ReadEndArray(); err != nil {
			t.Fatalf("ReadEndArray failed: %v", err)
		}

		if err := r.ReadEndArray(); err != nil {
			t.Fatalf("ReadEndArray failed: %v", err)
		}
	})
}

func TestWriteReadMap(t *testing.T) {
	t.Run("empty_map", func(t *testing.T) {
		w := NewCborWriter()
		if err := w.WriteStartMap(0); err != nil {
			t.Fatalf("WriteStartMap failed: %v", err)
		}
		if err := w.WriteEndMap(); err != nil {
			t.Fatalf("WriteEndMap failed: %v", err)
		}

		r := NewCborReader(w.Bytes())
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
	})

	t.Run("string_to_int_map", func(t *testing.T) {
		w := NewCborWriter()
		if err := w.WriteStartMap(2); err != nil {
			t.Fatalf("WriteStartMap failed: %v", err)
		}
		if err := w.WriteTextString("a"); err != nil {
			t.Fatalf("WriteTextString failed: %v", err)
		}
		if err := w.WriteInt64(1); err != nil {
			t.Fatalf("WriteInt64 failed: %v", err)
		}
		if err := w.WriteTextString("b"); err != nil {
			t.Fatalf("WriteTextString failed: %v", err)
		}
		if err := w.WriteInt64(2); err != nil {
			t.Fatalf("WriteInt64 failed: %v", err)
		}
		if err := w.WriteEndMap(); err != nil {
			t.Fatalf("WriteEndMap failed: %v", err)
		}

		r := NewCborReader(w.Bytes())
		length, err := r.ReadStartMap()
		if err != nil {
			t.Fatalf("ReadStartMap failed: %v", err)
		}
		if length != 2 {
			t.Errorf("got length %d, want 2", length)
		}

		key, err := r.ReadTextString()
		if err != nil {
			t.Fatalf("ReadTextString failed: %v", err)
		}
		if key != "a" {
			t.Errorf("got key %q, want 'a'", key)
		}
		val, err := r.ReadInt64()
		if err != nil {
			t.Fatalf("ReadInt64 failed: %v", err)
		}
		if val != 1 {
			t.Errorf("got value %d, want 1", val)
		}

		key, err = r.ReadTextString()
		if err != nil {
			t.Fatalf("ReadTextString failed: %v", err)
		}
		if key != "b" {
			t.Errorf("got key %q, want 'b'", key)
		}
		val, err = r.ReadInt64()
		if err != nil {
			t.Fatalf("ReadInt64 failed: %v", err)
		}
		if val != 2 {
			t.Errorf("got value %d, want 2", val)
		}

		if err := r.ReadEndMap(); err != nil {
			t.Fatalf("ReadEndMap failed: %v", err)
		}
	})
}

func TestWriteReadTag(t *testing.T) {
	tests := []struct {
		name string
		tag  CborTag
	}{
		{"datetime_string", TagDateTimeString},
		{"unix_time", TagUnixTime},
		{"unsigned_bignum", TagUnsignedBignum},
		{"uri", TagURI},
		{"self_described", TagSelfDescribedCbor},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewCborWriter()
			if err := w.WriteTag(tt.tag); err != nil {
				t.Fatalf("WriteTag failed: %v", err)
			}
			if err := w.WriteNull(); err != nil {
				t.Fatalf("WriteNull failed: %v", err)
			}

			r := NewCborReader(w.Bytes())
			tag, err := r.ReadTag()
			if err != nil {
				t.Fatalf("ReadTag failed: %v", err)
			}
			if tag != tt.tag {
				t.Errorf("got tag %d, want %d", tag, tt.tag)
			}
			if err := r.ReadNull(); err != nil {
				t.Fatalf("ReadNull failed: %v", err)
			}
		})
	}
}

func TestWriteReadBigInt(t *testing.T) {
	tests := []struct {
		name  string
		value *big.Int
	}{
		{"zero", big.NewInt(0)},
		{"positive", big.NewInt(12345)},
		{"negative", big.NewInt(-12345)},
		{"max_int64", big.NewInt(math.MaxInt64)},
		{"min_int64", big.NewInt(math.MinInt64)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewCborWriter()
			if err := w.WriteBigInt(tt.value); err != nil {
				t.Fatalf("WriteBigInt failed: %v", err)
			}

			r := NewCborReader(w.Bytes())
			got, err := r.ReadBigInt()
			if err != nil {
				t.Fatalf("ReadBigInt failed: %v", err)
			}
			if got.Cmp(tt.value) != 0 {
				t.Errorf("got %v, want %v", got, tt.value)
			}
		})
	}

	t.Run("very_large_positive", func(t *testing.T) {
		// Create a number larger than uint64
		value := new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil)

		w := NewCborWriter()
		if err := w.WriteBigInt(value); err != nil {
			t.Fatalf("WriteBigInt failed: %v", err)
		}

		r := NewCborReader(w.Bytes())
		got, err := r.ReadBigInt()
		if err != nil {
			t.Fatalf("ReadBigInt failed: %v", err)
		}
		if got.Cmp(value) != 0 {
			t.Errorf("got %v, want %v", got, value)
		}
	})

	t.Run("very_large_negative", func(t *testing.T) {
		// Create a number smaller than int64
		value := new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil)
		value.Neg(value)

		w := NewCborWriter()
		if err := w.WriteBigInt(value); err != nil {
			t.Fatalf("WriteBigInt failed: %v", err)
		}

		r := NewCborReader(w.Bytes())
		got, err := r.ReadBigInt()
		if err != nil {
			t.Fatalf("ReadBigInt failed: %v", err)
		}
		if got.Cmp(value) != 0 {
			t.Errorf("got %v, want %v", got, value)
		}
	})
}

func TestWriteReadDateTime(t *testing.T) {
	t.Run("datetime_string", func(t *testing.T) {
		original := time.Date(2024, 6, 15, 10, 30, 45, 0, time.UTC)

		w := NewCborWriter()
		if err := w.WriteDateTimeString(original); err != nil {
			t.Fatalf("WriteDateTimeString failed: %v", err)
		}

		r := NewCborReader(w.Bytes())
		got, err := r.ReadDateTimeString()
		if err != nil {
			t.Fatalf("ReadDateTimeString failed: %v", err)
		}
		if !got.Equal(original) {
			t.Errorf("got %v, want %v", got, original)
		}
	})

	t.Run("unix_time_integer", func(t *testing.T) {
		original := time.Unix(1718444445, 0)

		w := NewCborWriter()
		if err := w.WriteUnixTime(original); err != nil {
			t.Fatalf("WriteUnixTime failed: %v", err)
		}

		r := NewCborReader(w.Bytes())
		got, err := r.ReadUnixTime()
		if err != nil {
			t.Fatalf("ReadUnixTime failed: %v", err)
		}
		if !got.Equal(original) {
			t.Errorf("got %v, want %v", got, original)
		}
	})

	t.Run("unix_time_with_nanos", func(t *testing.T) {
		original := time.Unix(1718444445, 123456789)

		w := NewCborWriter()
		if err := w.WriteUnixTime(original); err != nil {
			t.Fatalf("WriteUnixTime failed: %v", err)
		}

		r := NewCborReader(w.Bytes())
		got, err := r.ReadUnixTime()
		if err != nil {
			t.Fatalf("ReadUnixTime failed: %v", err)
		}
		// Allow small differences due to float precision
		diff := got.Sub(original)
		if diff < -time.Microsecond || diff > time.Microsecond {
			t.Errorf("got %v, want %v (diff: %v)", got, original, diff)
		}
	})
}

func TestWriteReadUri(t *testing.T) {
	uri := "https://example.com/path?query=value"

	w := NewCborWriter()
	if err := w.WriteUri(uri); err != nil {
		t.Fatalf("WriteUri failed: %v", err)
	}

	r := NewCborReader(w.Bytes())
	tag, err := r.ReadTag()
	if err != nil {
		t.Fatalf("ReadTag failed: %v", err)
	}
	if tag != TagURI {
		t.Errorf("got tag %d, want %d", tag, TagURI)
	}
	got, err := r.ReadTextString()
	if err != nil {
		t.Fatalf("ReadTextString failed: %v", err)
	}
	if got != uri {
		t.Errorf("got %q, want %q", got, uri)
	}
}

func TestIndefiniteLengthArray(t *testing.T) {
	w := NewCborWriter()
	if err := w.WriteStartIndefiniteLengthArray(); err != nil {
		t.Fatalf("WriteStartIndefiniteLengthArray failed: %v", err)
	}
	if err := w.WriteInt64(1); err != nil {
		t.Fatalf("WriteInt64 failed: %v", err)
	}
	if err := w.WriteInt64(2); err != nil {
		t.Fatalf("WriteInt64 failed: %v", err)
	}
	if err := w.WriteInt64(3); err != nil {
		t.Fatalf("WriteInt64 failed: %v", err)
	}
	if err := w.WriteEndArray(); err != nil {
		t.Fatalf("WriteEndArray failed: %v", err)
	}

	r := NewCborReader(w.Bytes())
	length, err := r.ReadStartArray()
	if err != nil {
		t.Fatalf("ReadStartArray failed: %v", err)
	}
	if length != -1 {
		t.Errorf("expected indefinite length (-1), got %d", length)
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
}

func TestIndefiniteLengthMap(t *testing.T) {
	w := NewCborWriter()
	if err := w.WriteStartIndefiniteLengthMap(); err != nil {
		t.Fatalf("WriteStartIndefiniteLengthMap failed: %v", err)
	}
	if err := w.WriteTextString("key"); err != nil {
		t.Fatalf("WriteTextString failed: %v", err)
	}
	if err := w.WriteInt64(42); err != nil {
		t.Fatalf("WriteInt64 failed: %v", err)
	}
	if err := w.WriteEndMap(); err != nil {
		t.Fatalf("WriteEndMap failed: %v", err)
	}

	r := NewCborReader(w.Bytes())
	length, err := r.ReadStartMap()
	if err != nil {
		t.Fatalf("ReadStartMap failed: %v", err)
	}
	if length != -1 {
		t.Errorf("expected indefinite length (-1), got %d", length)
	}

	key, err := r.ReadTextString()
	if err != nil {
		t.Fatalf("ReadTextString failed: %v", err)
	}
	if key != "key" {
		t.Errorf("got key %q, want 'key'", key)
	}
	val, err := r.ReadInt64()
	if err != nil {
		t.Fatalf("ReadInt64 failed: %v", err)
	}
	if val != 42 {
		t.Errorf("got %d, want 42", val)
	}

	if err := r.ReadEndMap(); err != nil {
		t.Fatalf("ReadEndMap failed: %v", err)
	}
}

func TestIndefiniteLengthByteString(t *testing.T) {
	w := NewCborWriter()
	if err := w.WriteStartIndefiniteLengthByteString(); err != nil {
		t.Fatalf("WriteStartIndefiniteLengthByteString failed: %v", err)
	}
	if err := w.WriteByteStringChunk([]byte{1, 2, 3}); err != nil {
		t.Fatalf("WriteByteStringChunk failed: %v", err)
	}
	if err := w.WriteByteStringChunk([]byte{4, 5}); err != nil {
		t.Fatalf("WriteByteStringChunk failed: %v", err)
	}
	if err := w.WriteEndIndefiniteLengthByteString(); err != nil {
		t.Fatalf("WriteEndIndefiniteLengthByteString failed: %v", err)
	}

	r := NewCborReader(w.Bytes())
	got, err := r.ReadByteString()
	if err != nil {
		t.Fatalf("ReadByteString failed: %v", err)
	}
	expected := []byte{1, 2, 3, 4, 5}
	if !bytes.Equal(got, expected) {
		t.Errorf("got %v, want %v", got, expected)
	}
}

func TestIndefiniteLengthTextString(t *testing.T) {
	w := NewCborWriter()
	if err := w.WriteStartIndefiniteLengthTextString(); err != nil {
		t.Fatalf("WriteStartIndefiniteLengthTextString failed: %v", err)
	}
	if err := w.WriteTextStringChunk("Hello, "); err != nil {
		t.Fatalf("WriteTextStringChunk failed: %v", err)
	}
	if err := w.WriteTextStringChunk("World!"); err != nil {
		t.Fatalf("WriteTextStringChunk failed: %v", err)
	}
	if err := w.WriteEndIndefiniteLengthTextString(); err != nil {
		t.Fatalf("WriteEndIndefiniteLengthTextString failed: %v", err)
	}

	r := NewCborReader(w.Bytes())
	got, err := r.ReadTextString()
	if err != nil {
		t.Fatalf("ReadTextString failed: %v", err)
	}
	expected := "Hello, World!"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestSkipValue(t *testing.T) {
	w := NewCborWriter()
	if err := w.WriteStartArray(3); err != nil {
		t.Fatalf("WriteStartArray failed: %v", err)
	}
	if err := w.WriteInt64(1); err != nil {
		t.Fatalf("WriteInt64 failed: %v", err)
	}
	// Write a complex nested structure
	if err := w.WriteStartMap(1); err != nil {
		t.Fatalf("WriteStartMap failed: %v", err)
	}
	if err := w.WriteTextString("nested"); err != nil {
		t.Fatalf("WriteTextString failed: %v", err)
	}
	if err := w.WriteStartArray(2); err != nil {
		t.Fatalf("WriteStartArray failed: %v", err)
	}
	if err := w.WriteInt64(2); err != nil {
		t.Fatalf("WriteInt64 failed: %v", err)
	}
	if err := w.WriteInt64(3); err != nil {
		t.Fatalf("WriteInt64 failed: %v", err)
	}
	if err := w.WriteEndArray(); err != nil {
		t.Fatalf("WriteEndArray failed: %v", err)
	}
	if err := w.WriteEndMap(); err != nil {
		t.Fatalf("WriteEndMap failed: %v", err)
	}
	if err := w.WriteInt64(4); err != nil {
		t.Fatalf("WriteInt64 failed: %v", err)
	}
	if err := w.WriteEndArray(); err != nil {
		t.Fatalf("WriteEndArray failed: %v", err)
	}

	r := NewCborReader(w.Bytes())
	length, err := r.ReadStartArray()
	if err != nil {
		t.Fatalf("ReadStartArray failed: %v", err)
	}
	if length != 3 {
		t.Errorf("got length %d, want 3", length)
	}

	// Read first element
	val, err := r.ReadInt64()
	if err != nil {
		t.Fatalf("ReadInt64 failed: %v", err)
	}
	if val != 1 {
		t.Errorf("got %d, want 1", val)
	}

	// Skip the nested map
	if err := r.SkipValue(); err != nil {
		t.Fatalf("SkipValue failed: %v", err)
	}

	// Read last element
	val, err = r.ReadInt64()
	if err != nil {
		t.Fatalf("ReadInt64 failed: %v", err)
	}
	if val != 4 {
		t.Errorf("got %d, want 4", val)
	}

	if err := r.ReadEndArray(); err != nil {
		t.Fatalf("ReadEndArray failed: %v", err)
	}
}

func TestPeekState(t *testing.T) {
	w := NewCborWriter()
	if err := w.WriteInt64(42); err != nil {
		t.Fatalf("WriteInt64 failed: %v", err)
	}

	r := NewCborReader(w.Bytes())

	// Peek multiple times
	for i := 0; i < 3; i++ {
		state, err := r.PeekState()
		if err != nil {
			t.Fatalf("PeekState failed: %v", err)
		}
		if state != StateUnsignedInteger {
			t.Errorf("got state %v, want %v", state, StateUnsignedInteger)
		}
	}

	// Now read
	val, err := r.ReadInt64()
	if err != nil {
		t.Fatalf("ReadInt64 failed: %v", err)
	}
	if val != 42 {
		t.Errorf("got %d, want 42", val)
	}

	// State should now be finished
	state, err := r.PeekState()
	if err != nil {
		t.Fatalf("PeekState failed: %v", err)
	}
	if state != StateFinished {
		t.Errorf("got state %v, want %v", state, StateFinished)
	}
}

func TestSimpleValue(t *testing.T) {
	tests := []struct {
		name  string
		value SimpleValue
	}{
		{"value_16", SimpleValue(16)},
		{"value_32", SimpleValue(32)},
		{"value_255", SimpleValue(255)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewCborWriter()
			if err := w.WriteSimpleValue(tt.value); err != nil {
				t.Fatalf("WriteSimpleValue failed: %v", err)
			}

			r := NewCborReader(w.Bytes())
			got, err := r.ReadSimpleValue()
			if err != nil {
				t.Fatalf("ReadSimpleValue failed: %v", err)
			}
			if got != tt.value {
				t.Errorf("got %d, want %d", got, tt.value)
			}
		})
	}
}

func TestTryReadNull(t *testing.T) {
	t.Run("is_null", func(t *testing.T) {
		w := NewCborWriter()
		if err := w.WriteNull(); err != nil {
			t.Fatalf("WriteNull failed: %v", err)
		}

		r := NewCborReader(w.Bytes())
		isNull, err := r.TryReadNull()
		if err != nil {
			t.Fatalf("TryReadNull failed: %v", err)
		}
		if !isNull {
			t.Errorf("expected true, got false")
		}
	})

	t.Run("is_not_null", func(t *testing.T) {
		w := NewCborWriter()
		if err := w.WriteInt64(42); err != nil {
			t.Fatalf("WriteInt64 failed: %v", err)
		}

		r := NewCborReader(w.Bytes())
		isNull, err := r.TryReadNull()
		if err != nil {
			t.Fatalf("TryReadNull failed: %v", err)
		}
		if isNull {
			t.Errorf("expected false, got true")
		}
		// Should still be able to read the value
		val, err := r.ReadInt64()
		if err != nil {
			t.Fatalf("ReadInt64 failed: %v", err)
		}
		if val != 42 {
			t.Errorf("got %d, want 42", val)
		}
	})
}

func TestCanonicalModeRejectsIndefiniteLength(t *testing.T) {
	w := NewCborWriter(WithConformanceMode(ConformanceCanonical))

	err := w.WriteStartIndefiniteLengthArray()
	if err != ErrIndefiniteLengthNotAllowed {
		t.Errorf("expected ErrIndefiniteLengthNotAllowed, got %v", err)
	}

	err = w.WriteStartIndefiniteLengthMap()
	if err != ErrIndefiniteLengthNotAllowed {
		t.Errorf("expected ErrIndefiniteLengthNotAllowed, got %v", err)
	}

	err = w.WriteStartIndefiniteLengthByteString()
	if err != ErrIndefiniteLengthNotAllowed {
		t.Errorf("expected ErrIndefiniteLengthNotAllowed, got %v", err)
	}

	err = w.WriteStartIndefiniteLengthTextString()
	if err != ErrIndefiniteLengthNotAllowed {
		t.Errorf("expected ErrIndefiniteLengthNotAllowed, got %v", err)
	}
}

func TestNestingDepthLimit(t *testing.T) {
	w := NewCborWriter(WithMaxNestingDepth(3))

	if err := w.WriteStartArray(1); err != nil {
		t.Fatalf("WriteStartArray 1 failed: %v", err)
	}
	if err := w.WriteStartArray(1); err != nil {
		t.Fatalf("WriteStartArray 2 failed: %v", err)
	}
	if err := w.WriteStartArray(1); err != nil {
		t.Fatalf("WriteStartArray 3 failed: %v", err)
	}

	// This should fail
	err := w.WriteStartArray(1)
	if err != ErrNestingDepthExceeded {
		t.Errorf("expected ErrNestingDepthExceeded, got %v", err)
	}
}

func TestReadEncodedValue(t *testing.T) {
	w := NewCborWriter()
	if err := w.WriteStartArray(2); err != nil {
		t.Fatalf("WriteStartArray failed: %v", err)
	}
	if err := w.WriteInt64(1); err != nil {
		t.Fatalf("WriteInt64 failed: %v", err)
	}
	if err := w.WriteInt64(2); err != nil {
		t.Fatalf("WriteInt64 failed: %v", err)
	}
	if err := w.WriteEndArray(); err != nil {
		t.Fatalf("WriteEndArray failed: %v", err)
	}

	original := w.BytesCopy()

	r := NewCborReader(original)
	encoded, err := r.ReadEncodedValue()
	if err != nil {
		t.Fatalf("ReadEncodedValue failed: %v", err)
	}
	if !bytes.Equal(encoded, original) {
		t.Errorf("encoded value doesn't match original")
	}
}

func TestResetWriter(t *testing.T) {
	w := NewCborWriter()
	if err := w.WriteInt64(42); err != nil {
		t.Fatalf("WriteInt64 failed: %v", err)
	}

	first := w.BytesCopy()

	w.Reset()
	if err := w.WriteInt64(123); err != nil {
		t.Fatalf("WriteInt64 failed: %v", err)
	}

	second := w.BytesCopy()

	if bytes.Equal(first, second) {
		t.Errorf("expected different results after reset")
	}

	r := NewCborReader(second)
	val, err := r.ReadInt64()
	if err != nil {
		t.Fatalf("ReadInt64 failed: %v", err)
	}
	if val != 123 {
		t.Errorf("got %d, want 123", val)
	}
}

func TestResetReader(t *testing.T) {
	w := NewCborWriter()
	if err := w.WriteInt64(42); err != nil {
		t.Fatalf("WriteInt64 failed: %v", err)
	}

	r := NewCborReader(w.Bytes())
	val, err := r.ReadInt64()
	if err != nil {
		t.Fatalf("ReadInt64 failed: %v", err)
	}
	if val != 42 {
		t.Errorf("got %d, want 42", val)
	}

	r.Reset()
	val, err = r.ReadInt64()
	if err != nil {
		t.Fatalf("ReadInt64 after reset failed: %v", err)
	}
	if val != 42 {
		t.Errorf("got %d, want 42", val)
	}
}
