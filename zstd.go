package cache

import (
	"github.com/klauspost/compress/zstd"
)

// ZSTDTranscoder provides transparent Zstandard compression and decompression using pre-initialized
// encoder and decoder instances. It is designed for high-performance scenarios where the same
// transcoder instance is reused across many operations. The implementation maintains internal
// state (the encoder and decoder) and is safe for concurrent use because the underlying zstd
// library guarantees thread-safety of its Writer and Reader types when properly configured.
type ZSTDTranscoder struct {
	Encoder *zstd.Encoder
	Decoder *zstd.Decoder
}

// NewZSTDTranscoder creates a new ZSTDTranscoder with a fully configured encoder and decoder.
// The encoder is initialized with the highest compression level for optimal ratio.
// Both objects are created with nil writers/readers because EncodeAll/DecodeAll do not require
// an io.Writer/io.Reader – they operate on complete byte slices.
// On any error during initialization, resources are cleaned up and the error is propagated.
func NewZSTDTranscoder() (*ZSTDTranscoder, error) {
	enc, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
	if err != nil {
		return nil, err
	}

	dec, err := zstd.NewReader(nil)
	if err != nil {
		_ = enc.Close()
		return nil, err
	}

	return &ZSTDTranscoder{Encoder: enc, Decoder: dec}, nil
}

// Compress accepts arbitrary input data and returns its fully Zstandard-compressed representation.
// It uses the fast EncodeAll path that compresses the entire input in one call and appends
// the result to a freshly allocated destination buffer of appropriate capacity.
// The returned slice is always a new allocation owned by the caller.
func (t *ZSTDTranscoder) Compress(src []byte) ([]byte, error) {
	return t.Encoder.EncodeAll(src, make([]byte, 0, len(src))), nil
}

// Decompress accepts Zstandard-compressed data and returns the original uncompressed bytes.
// It uses the convenient DecodeAll method which handles the complete decompression in a single
// operation. The destination buffer is managed internally; the returned slice is a new allocation
// owned by the caller. Any error during decompression (corrupted data, incomplete input, etc.)
// is returned to the caller.
func (t *ZSTDTranscoder) Decompress(src []byte) ([]byte, error) {
	return t.Decoder.DecodeAll(src, nil)
}
