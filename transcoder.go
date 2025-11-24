package cache

import "errors"

type Transcoder[T any] interface {
	Encode(T) (string, error)
	Decode(string) (T, error)

	Close()
}

type PipelineTranscoder[T any] struct {
	jsonTranscoder   *JSONTranscoder[T]
	zstdTranscoder   *ZSTDTranscoder
	binaryTranscoder *Base64Transcoder
}

func NewPipelineTranscoder[T any]() *PipelineTranscoder[T] {
	zstdTranscoder, _ := NewZSTDTranscoder()
	return &PipelineTranscoder[T]{
		jsonTranscoder:   NewJSONTranscoder[T](),
		zstdTranscoder:   zstdTranscoder,
		binaryTranscoder: NewBase64Transcoder(),
	}
}

func (t *PipelineTranscoder[T]) Encode(src T) (string, error) {
	jsonBytes, err := t.jsonTranscoder.Marshal(src)
	if err != nil {
		return "", errors.Join(errors.New("failed to marshal JSON"), err)
	}

	compressedBytes, err := t.zstdTranscoder.Compress(jsonBytes)
	if err != nil {
		return "", errors.Join(errors.New("failed to compress with Zstd"), err)
	}

	return t.binaryTranscoder.Encode(compressedBytes)
}

func (t *PipelineTranscoder[T]) Decode(src string) (T, error) {
	var entry T

	compressedBytes, err := t.binaryTranscoder.Decode(src)
	if err != nil {
		return entry, errors.Join(errors.New("failed to decode Base64"), err)
	}

	jsonBytes, err := t.zstdTranscoder.Decompress(compressedBytes)
	if err != nil {
		return entry, errors.Join(errors.New("failed to decompress Zstd"), err)
	}

	return t.jsonTranscoder.Unmarshal(jsonBytes)
}

func (t *PipelineTranscoder[T]) Close() {
	t.zstdTranscoder.Close()
}
