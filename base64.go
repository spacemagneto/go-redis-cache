package cache

import "encoding/base64"

// Base64Transcoder provides a straightforward implementation of Base64 encoding and decoding
// using Go's standard library encoding/base64 package with the standard padding rules (RFC 4648).
type Base64Transcoder struct{}

// NewBase64Transcoder creates and returns a new instance of Base64Transcoder.
// The returned object has no internal state and can be safely shared across the application.
// It is provided as a constructor to maintain a consistent creation pattern with other transcoder types.
func NewBase64Transcoder() *Base64Transcoder {
	return &Base64Transcoder{}
}

// Encode converts the given byte slice into a Base64-encoded string using the standard encoding.
// The result includes padding characters (=) when necessary to comply with RFC 4648.
// No error is ever returned because the standard encoding is guaranteed to succeed for any input
func (b *Base64Transcoder) Encode(src []byte) (string, error) {
	return base64.StdEncoding.EncodeToString(src), nil
}

// Decode converts a Base64-encoded string back into its original byte representation.
// It accepts both padded and un-padded input (the standard decoder is tolerant of missing padding).
// If the input contains invalid characters or incorrect padding, a non-nil error is returned.
func (b *Base64Transcoder) Decode(src string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(src)
}
