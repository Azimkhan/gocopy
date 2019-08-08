package main

import (
	"gotest.tools/assert"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestCopy(t *testing.T) {
	// === Copy without offset
	copyErr := Copy("src", "dest", 2048, 0)
	assert.NilError(t, copyErr)
	read, _ := ioutil.ReadFile("dest")
	shouldBe, _ := ioutil.ReadFile("r1")
	assert.Equal(t, true, reflect.DeepEqual(shouldBe, read))

	// === Copy with offset
	copyErr = Copy("src", "dest", 2048, 3123123)
	assert.NilError(t, copyErr)
	read, _ = ioutil.ReadFile("dest")
	shouldBe, _ = ioutil.ReadFile("r2")
	assert.Equal(t, true, reflect.DeepEqual(shouldBe, read))
}
