package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var cryptionTests = []struct {
	in string
}{
	{"asdf"},
	{"relay"},
	{"12309fasdfa"},
	{")(*&^%$##@@!{}[],.<>??/\\"},
}

func TestEncryptDecrypt(t *testing.T) {
	for _, tt := range cryptionTests {
		encrypted, err := Encrypt(tt.in)
		assert.Nil(t, err)
		decrypted, err := Decrypt(encrypted)
		assert.Nil(t, err)
		assert.Equal(t, tt.in, decrypted)
	}
}
