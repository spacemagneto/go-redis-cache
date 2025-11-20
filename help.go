package cache

import "time"

// sample payloads of different sizes and compressibility
var (
	// Small payload – typical short event or config object
	smallPayload = []byte(`{"id":12345,"user":"alice","action":"login","timestamp":"2025-08-09T12:34:56Z","metadata":{"ip":"203.0.113.42"}}`)

	// Medium payload – typical log entry or message with nested data
	mediumPayload = []byte(`{
		"event":"order.created",
		"order_id":"ord_987654321",
		"user_id":556677,
		"items":[{"sku":"ABC-123","qty":2,"price":49.99},{"sku":"XYZ-789","qty":1,"price":299.0}],
		"total":398.98,
		"currency":"USD",
		"timestamp":"2025-08-09T14:22:18.123Z",
		"metadata":{"session":"sess_abc123","referrer":"https://example.com/checkout"}
	}`)

	// Large payload – highly compressible repeated text (worst-case best compression)
	largeCompressiblePayload = []byte(makeRepeatedString(
		`{"level":"info","ts":"2025-08-09T10:00:00.000Z","msg":"request completed","path":"/api/v1/users","method":"GET","duration_ms":12} `,
		1000,
	))

	// Large random-like payload – simulates encrypted data or already-compressed content (worst-case ratio)
	largeRandomPayload = func() []byte {
		b := make([]byte, 500_000)
		for i := range b {
			b[i] = byte(time.Now().UnixNano() >> (i % 8))
		}
		return b
	}()
)

// helper to generate repeated compressible text
func makeRepeatedString(s string, n int) string {
	result := make([]byte, 0, len(s)*n)
	for i := 0; i < n; i++ {
		result = append(result, s...)
	}
	return string(result)
}
