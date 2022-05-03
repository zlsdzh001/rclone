package readers

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Check interface
var _ io.ReadSeeker = &FakeSeeker{}

func TestFakeSeeker(t *testing.T) {
	in := bytes.NewBufferString("hello")
	buf := make([]byte, 16)
	r := NewFakeSeeker(in, 5)

	// check the seek offset is as passed in
	checkPos := func(pos int64) {
		abs, err := r.Seek(0, io.SeekCurrent)
		require.NoError(t, err)
		assert.Equal(t, pos, abs)
	}

	// Test some seeking
	checkPos(0)

	abs, err := r.Seek(2, io.SeekStart)
	require.NoError(t, err)
	assert.Equal(t, int64(2), abs)
	checkPos(2)

	abs, err = r.Seek(-1, io.SeekEnd)
	require.NoError(t, err)
	assert.Equal(t, int64(4), abs)
	checkPos(4)

	// Check can't read if not at start
	_, err = r.Read(buf)
	require.ErrorContains(t, err, "not at start")

	// Seek back to start
	abs, err = r.Seek(-4, io.SeekCurrent)
	require.NoError(t, err)
	assert.Equal(t, int64(0), abs)
	checkPos(0)

	_, err = r.Seek(42, 17)
	require.ErrorContains(t, err, "invalid whence")

	_, err = r.Seek(-1, io.SeekStart)
	require.ErrorContains(t, err, "negative position")

	// Test reading now seeked back to the start
	n, err := r.Read(buf)
	require.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, []byte("hello"), buf[:5])

	// Seeking should give an error now
	_, err = r.Seek(-1, io.SeekEnd)
	require.ErrorContains(t, err, "after reading")
}
