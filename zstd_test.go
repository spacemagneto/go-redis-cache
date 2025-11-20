package cache

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestZSTDTranscoder is the table-driven test for the ZSTDTranscoder type.
// It ensures that a single transcoder instance can be safely reused for multiple
// compression and decompression operations, that round-trip fidelity is preserved
// for various input sizes including empty data, and that decompression correctly
// fails with an error when given corrupted or truncated compressed data.
// The test creates the transcoder once before running any cases to reflect the
// intended high-performance usage pattern where initialization cost is amortized
// over many operations.
func TestZSTDTranscoder(t *testing.T) {
	t.Parallel()

	transcoder, err := NewZSTDTranscoder()
	assert.NoError(t, err, "Failed init new zstd transcoder")

	defer transcoder.encoder.Close()
	defer transcoder.decoder.Close()

	cases := []struct {
		name        string
		input       []byte
		corruptFunc func([]byte) []byte
		wantErr     bool
	}{
		{name: "Empty input", input: []byte{}},
		{name: "Short text", input: []byte("hello world")},
		{name: "Long repeatable data", input: []byte(strings.Repeat("grok-xai-2025-", 100))},
		{name: "Single byte", input: []byte{0x42}},
		{name: "Corrupted data - bit flip", input: []byte("robustness test"),
			corruptFunc: func(b []byte) []byte {
				if len(b) > 10 {
					c := append([]byte(nil), b...)
					c[10] ^= 0xff
					return c
				}
				return b
			},
			wantErr: true,
		},
		{name: "Truncated compressed data", input: []byte("this will not survive truncation"),
			corruptFunc: func(b []byte) []byte {
				if len(b) > 8 {
					return b[:8]
				}
				return b
			},
			wantErr: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			compressed, err := transcoder.Compress(tt.input)
			assert.NoError(t, err, "Compress must succeed for valid input")

			dataToDecompress := compressed
			if tt.corruptFunc != nil {
				dataToDecompress = tt.corruptFunc(compressed)
			}

			decompressed, err := transcoder.Decompress(dataToDecompress)

			if tt.wantErr {
				assert.Error(t, err, "Decompress must return an error for corrupted or invalid input")
				return
			}

			// For successful cases, verify full round-trip fidelity.
			// The zstd library returns []byte(nil) for empty input while we pass []byte{}.
			// Convert nil → empty slice so assert.Equal works perfectly in all cases.
			if decompressed == nil {
				decompressed = []byte{}
			}

			assert.NoError(t, err, "Decompress must succeed for valid compressed data")
			assert.Equal(t, tt.input, decompressed, "Failed: decompressed data does not match original input")
		})
	}
}
