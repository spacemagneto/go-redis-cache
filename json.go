package cache

import "github.com/goccy/go-json"

// JSONTranscoder provides a generic, stateless wrapper for JSON serialization and deserialization.
// It enables type-safe marshalling and unmarshalling of arbitrary Go values while preserving
// full compatibility with standard JSON encoding rules (struct tags, custom Marshal/Unmarshal, etc.).
type JSONTranscoder[T any] struct{}

// NewJSONTranscoder creates a new instance of JSONTranscoder for the specified type T.
// Because the implementation has no internal state, the returned instance can be safely
// shared and reused throughout the application. The constructor exists to provide a
// consistent instantiation pattern with other transcoder implementations.
func NewJSONTranscoder[T any]() *JSONTranscoder[T] {
	return &JSONTranscoder[T]{}
}

// Marshal converts the given value of type T into its JSON byte representation.
// The operation respects all standard JSON marshaling semantics and produces valid UTF-8 output.
func (t *JSONTranscoder[T]) Marshal(src T) ([]byte, error) {
	return json.Marshal(src)
}

// Unmarshal parses JSON data from the provided byte slice and populates a new value of type T.
// On success it returns the decoded value, on failure it returns the zero value of T together
// with the parsing or assignment error.
func (t *JSONTranscoder[T]) Unmarshal(src []byte) (T, error) {
	var entry T
	err := json.Unmarshal(src, &entry)
	return entry, err
}
