package cache

import (
	"math/rand"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestIsNil(t *testing.T) {
	t.Parallel()

	t.Run("NilValuesOfNilKinds", func(t *testing.T) {
		assert.True(t, IsNil[*int](nil))
		assert.True(t, IsNil[[]string](nil))
		assert.True(t, IsNil[map[int]bool](nil))
		assert.True(t, IsNil[chan struct{}](nil))
		assert.True(t, IsNil[func()](nil))
		assert.True(t, IsNil[any](nil))
		assert.True(t, IsNil[error](nil))
		assert.False(t, IsNil[unsafe.Pointer](nil))
	})

	t.Run("NonNilValuesOfNilKinds", func(t *testing.T) {
		x := rand.Int()
		assert.False(t, IsNil(&x))
		assert.False(t, IsNil([]int{1, 2}))
		assert.False(t, IsNil([]int{}))
		assert.False(t, IsNil(map[string]int{}))
		assert.False(t, IsNil(make(chan int)))
		assert.False(t, IsNil(func() {}))
		assert.False(t, IsNil[any]("hello"))
	})

	t.Run("ZeroValuesOfNonNilKinds", func(t *testing.T) {
		assert.False(t, IsNil(0))
		assert.False(t, IsNil(3.14))
		assert.False(t, IsNil(""))
		assert.False(t, IsNil(true))
		assert.False(t, IsNil(false))

		type User struct{ Name string }
		assert.False(t, IsNil(User{}))
		assert.False(t, IsNil([3]int{1, 2, 3}))
	})

	t.Run("TypedNilPointersAndInterfaces", func(t *testing.T) {
		var p *int = nil
		assert.True(t, IsNil(p))

		var i any = (*string)(nil)
		assert.True(t, IsNil(i))

		var err error = nil
		assert.True(t, IsNil(err))
	})

	t.Run("InterfaceHoldingConcreteValue", func(t *testing.T) {
		var i any = "not nil"
		assert.False(t, IsNil(i))

		var j any = 100
		assert.False(t, IsNil(j))

		s := "hello"
		var k any = &s
		assert.False(t, IsNil(k))
	})

	t.Run("NilInterfaceVsNilPointerInsideInterface", func(t *testing.T) {
		var nilInterface any = nil
		assert.True(t, IsNil(nilInterface))

		var nilPtrInInterface any = (*int)(nil)
		assert.True(t, IsNil(nilPtrInInterface))

		x := rand.Int()
		var nonNilInInterface any = &x
		assert.False(t, IsNil(nonNilInInterface))
	})
}
