package ringio

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadWrites(t *testing.T) {
	b := New(9)
	n, err := b.Write([]byte("hi"))
	assert.NoError(t, err)
	assert.Equal(t, 2, n)

	n, err = b.Write([]byte("this will hit the capacity of the buffer"))
	assert.Error(t, err)
	assert.Equal(t, io.ErrShortWrite, err)
	assert.Equal(t, 7, n)

	buf := make([]byte, 9)
	n, err = b.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, []byte("hithis wi"), buf)
}
