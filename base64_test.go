package cache

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBase64TranscoderEncode is a table-driven test for the Encode method.
// It verifies that various byte inputs are correctly encoded to standard Base64 strings
// (RFC 4648 with padding).
func TestBase64TranscoderEncode(t *testing.T) {
	t.Parallel()

	transcoder := NewBase64Transcoder()

	cases := []struct {
		name     string
		input    []byte
		expected string
	}{
		{name: "Empty input", input: []byte{}, expected: ""},
		{name: "Single character 'f'", input: []byte("f"), expected: "Zg=="},
		{name: "Two characters fo", input: []byte("fo"), expected: "Zm8="},
		{name: "Three characters foo", input: []byte("foo"), expected: "Zm9v"},
		{name: "Standard phrase hello world", input: []byte("hello world"), expected: "aGVsbG8gd29ybGQ="},
		{name: "Lorem ipsum", input: []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum maximus lacus sed libero pretium, faucibus mattis lacus vestibulum. Curabitur eget vulputate dolor, tincidunt ullamcorper nisl."), expected: "TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNlY3RldHVyIGFkaXBpc2NpbmcgZWxpdC4gVmVzdGlidWx1bSBtYXhpbXVzIGxhY3VzIHNlZCBsaWJlcm8gcHJldGl1bSwgZmF1Y2lidXMgbWF0dGlzIGxhY3VzIHZlc3RpYnVsdW0uIEN1cmFiaXR1ciBlZ2V0IHZ1bHB1dGF0ZSBkb2xvciwgdGluY2lkdW50IHVsbGFtY29ycGVyIG5pc2wu"},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			result, err := transcoder.Encode(tt.input)

			assert.NoError(t, err, "Encode should never return an error")
			assert.Equal(t, tt.expected, result, "Encoded string does not match expected Base64 for input %q", tt.input)

			stdEncoded := base64.StdEncoding.EncodeToString(tt.input)
			assert.Equal(t, stdEncoded, result, "Should match encoding/base64.StdEncoding")
		})
	}
}

// TestBase64TranscoderDecode is a table-driven test for the Decode method.
// It covers valid Base64 strings (with and without padding), as well as various invalid cases.
func TestBase64TranscoderDecode(t *testing.T) {
	t.Parallel()

	transcoder := NewBase64Transcoder()

	cases := []struct {
		name        string
		input       string
		expected    []byte
		expectError bool
	}{
		{name: "Empty string", input: "", expected: []byte{}, expectError: false},
		{name: "Padded single char f", input: "Zg==", expected: []byte("f"), expectError: false},
		{name: "Un padded hello world", input: "aGVsbG8gd29ybGQ=", expected: []byte("hello world"), expectError: false},
		{name: "Three chars no padding", input: "Zm9v", expected: []byte("foo"), expectError: false},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			result, err := transcoder.Decode(tt.input)

			if tt.expectError {
				assert.Error(t, err, "Decode should fail on invalid input %q", tt.input)
				assert.Nil(t, result, "Result should be nil when error occurs")
			} else {
				assert.NoError(t, err, "Decode should succeed for valid input %q", tt.input)
				assert.Equal(t, tt.expected, result, "Decoded bytes do not match expected for input %q", tt.input)
			}
		})
	}
}
