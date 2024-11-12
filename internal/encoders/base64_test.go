package encoders

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_EncodeBase64(t *testing.T) {
	str := "user:W07QCRPA4"
	encodedStr := EncodeToBase64(str)
	assert.EqualValues(t, encodedStr, "dXNlcjpXMDdRQ1JQQTQ=")
}


func Test_DecodeBase64(t *testing.T) {
	str := "dXNlcjpXMDdRQ1JQQTQ="
	encodedStr, err := DecodeBase64(str)
	assert.EqualValues(t, encodedStr, "user:W07QCRPA4")
	assert.Nil(t, err)
}
