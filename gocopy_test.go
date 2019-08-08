package main

import (
	"gotest.tools/assert"
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	var fileErr error
	var copyErr error
	var n int
	var file *os.File
	buf := make([]byte, 1)

	// === Copy without offset
	copyErr = Copy("src", "dest", 2048, 0)
	assert.NilError(t, copyErr)
	file, fileErr = os.Open("dest")
	assert.NilError(t, fileErr)

	// check first and last byte of a file
	n, fileErr = file.ReadAt(buf, int64(2047))
	assert.NilError(t, fileErr)
	assert.Equal(t, 1, n)
	assert.Equal(t, uint8(10), buf[0]) // should be LF

	n, fileErr = file.ReadAt(buf, int64(0))
	assert.NilError(t, fileErr)
	assert.Equal(t, 1, n)
	assert.Equal(t, uint8(97), buf[0]) // should be "a"

	// === Copy with offset
	copyErr = Copy("src", "dest", 2048, 1024)
	assert.NilError(t, copyErr)
	file, fileErr = os.Open("dest")
	assert.NilError(t, fileErr)

	// check first and last byte of a file
	n, fileErr = file.ReadAt(buf, int64(1023))
	assert.NilError(t, fileErr)
	assert.Equal(t, 1, n)
	assert.Equal(t, uint8(10), buf[0]) // should be LF

	n, fileErr = file.ReadAt(buf, int64(0))
	assert.NilError(t, fileErr)
	assert.Equal(t, 1, n)
	assert.Equal(t, uint8(99), buf[0]) // should be "c"
}
