package ringio

import (
	"io"
	"testing"

	"github.com/antifuchs/o"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadBoundedWrites(t *testing.T) {
	t.Parallel()

	b := New(9, false)
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

func TestReadOverwrites(t *testing.T) {
	t.Parallel()

	b := New(9, true)
	n, err := b.Write([]byte("hi"))
	assert.NoError(t, err)
	assert.Equal(t, 2, n)

	n, err = b.Write([]byte("this will hit the capacity of the buffer"))
	assert.NoError(t, err)
	assert.Equal(t, 40, n)

	buf := make([]byte, 9)
	n, err = b.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, []byte("he buffer"), buf)
}
