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

func TestParallel(t *testing.T) {
	t.Parallel()

	b := New(27, false)
	quit := make(chan struct{})
	write := func(toWrite []byte) {
		for {
			select {
			case <-quit:
				return
			default:
				b.Write(toWrite)
			}
		}
	}
	go write([]byte("abc"))
	go write([]byte("abc"))

	for i := 0; i < 1000; i++ {
		didRead := make([]byte, 6)
		n, err := b.Read(didRead)
		if err == o.ErrEmpty {
			continue
		}
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		switch n {
		case 3:
			assert.Equal(t, []byte("abc"), didRead[0:3])
		case 6:
			assert.Equal(t, []byte("abcabc"), didRead)
		default:
			t.Fatalf("Read %d bytes, expected 3 or 6", n)
		}
	}

	close(quit)
}

func TestReset(t *testing.T) {
	t.Parallel()

	b := New(8, true)
	n, err := b.Write([]byte("hi this is a test"))
	require.NoError(t, err)
	assert.Equal(t, 17, n)

	read := make([]byte, 4)
	n, err = b.Read(read)
	require.NoError(t, err)
	assert.Equal(t, 4, n)
	b.Reset()

	n, err = b.Read(read)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 0, n)
}
