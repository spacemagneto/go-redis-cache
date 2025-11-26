package cache

type Transcoder[T any] interface {
	Encode(T) (string, error)

	Decode(string) (T, error)
}
