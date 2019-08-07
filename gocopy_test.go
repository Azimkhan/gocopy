package main

import (
	"gotest.tools/assert"
	"testing"
)

func TestCopy(t *testing.T) {
	var err error
	err = Copy("src", "dest", 2048, 0)
	assert.NilError(t, err)

	err = Copy("src", "dest", 2048, 1024)
	assert.NilError(t, err)
}
