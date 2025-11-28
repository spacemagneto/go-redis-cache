package cache

import (
	"github.com/goccy/go-json"
)

type Transcoder[T any] interface {
	Encode(T) (string, error)

	Decode(string) (T, error)
}

type defaultTranscoder[T any] struct{}

func (defaultTranscoder[T]) Encode(src T) (string, error) {
	var bytes []byte
	var err error

	bytes, err = json.Marshal(src)

	return string(bytes), err
}

func (defaultTranscoder[T]) Decode(src string) (T, error) {
	var entry T
	var err error

	if err = json.Unmarshal([]byte(src), &entry); err != nil {
		return entry, err
	}

	return entry, nil
}
