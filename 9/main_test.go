package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnpack(t *testing.T) {
	s1, err1 := Unpack("a4bc2d5e")
	assert.Equal(t, "aaaabccddddde", s1)
	assert.NoError(t, err1)

	s2, err2 := Unpack("abcd")
	assert.Equal(t, "abcd", s2)
	assert.NoError(t, err2)

	s3, err3 := Unpack("45")
	assert.Empty(t, s3)
	assert.Error(t, err3)

	s4, err4 := Unpack("")
	assert.Equal(t, "", s4)
	assert.NoError(t, err4)

	s5, err5 := Unpack("qwe\\4\\5")
	assert.Equal(t, "qwe45", s5)
	assert.NoError(t, err5)

	s6, err6 := Unpack("qwe\\45")
	assert.Equal(t, "qwe44444", s6)
	assert.NoError(t, err6)

	s7, err7 := Unpack("4wa2")
	assert.Empty(t, s7)
	assert.Error(t, err7)

	s8, err8 := Unpack("\\35")
	assert.Equal(t, "33333", s8)
	assert.NoError(t, err8)
}
