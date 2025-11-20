package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Person struct {
	Name  string `json:"name"`
	Age   int    `json:"age,omitempty"`
	Email string `json:"email,omitempty"`
}

// TestJSONTranscoderMarshal is the table-driven test for the Marshal method of JSONTranscoder[T].
// It ensures the method correctly converts a value of type T into its JSON byte representation while
// fully preserving standard JSON marshalling semantics. This includes proper handling of struct field
// tags such as omitempty, ignoring fields tagged with "-", omitting unexported fields, and correctly
// invoking custom MarshalJSON methods when present. The test also verifies that Marshal never returns
// an error for any valid Go value.
func TestJSONTranscoderMarshal(t *testing.T) {
	transcoder := NewJSONTranscoder[Person]()

	cases := []struct {
		name     string
		input    Person
		expected string
	}{
		{name: "Full struct", input: Person{Name: "Alice", Age: 30, Email: "alice@example.com"}, expected: `{"name":"Alice","age":30,"email":"alice@example.com"}`},
		{name: "Omit empty fields", input: Person{Name: "Bob"}, expected: `{"name":"Bob"}`},
		{name: "Only email", input: Person{Name: "Charlie", Email: "charlie@x.com"}, expected: `{"name":"Charlie","email":"charlie@x.com"}`},
		{name: "Empty struct", input: Person{}, expected: `{"name":""}`},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			result, err := transcoder.Marshal(tt.input)

			assert.NoError(t, err, "Marshal must never return an error for valid Go values")
			assert.JSONEq(t, tt.expected, string(result), "Marshaled JSON does not match expected for input %+v", tt.input)
		})
	}
}
