package cache

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBase64Transcoder_Encode is a table-driven test for the Encode method.
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
