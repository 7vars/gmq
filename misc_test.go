package gmq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLenght(t *testing.T) {
	var l uint64 = 301
	b := encodeLength(l)
	assert.Equal(t, 2, len(b))
	assert.Equal(t, []byte{173, 2}, b)

	lx, buf := extractLength(append(b, 8, 9, 10))
	assert.Equal(t, l, lx)
	assert.Equal(t, []byte{8, 9, 10}, buf)
}

func TestString(t *testing.T) {
	s := `Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem`

	b := encodeString(s)
	assert.Equal(t, byte(1), b[0])
	assert.Equal(t, byte(45), b[1])

	res, buf := extractString(append(b, 8, 9, 10))
	assert.Equal(t, s, res)
	assert.Equal(t, []byte{8, 9, 10}, buf)

	s = `Hello 世界`
	b = encodeString(s)
	res, buf = extractString(append(b, 8, 9, 10))
	assert.Equal(t, s, res)
	assert.Equal(t, []byte{8, 9, 10}, buf)
}

func TestUint16(t *testing.T) {
	var i uint16 = 30
	b := encodeUint16(i)
	assert.Equal(t, 2, len(b))
	assert.Equal(t, []byte{0, 30}, b)

	ix, buf := extractUint16(append(b, 8, 9, 10))
	assert.Equal(t, i, ix)
	assert.Equal(t, []byte{8, 9, 10}, buf)
}
