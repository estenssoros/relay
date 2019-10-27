package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	config, err := Load()
	assert.Nil(t, err)
	assert.NotEmpty(t, config)
	fmt.Println(config)
}

func TestCipherKeyBytes(t *testing.T) {
	b, err := CipherKeyBytes()
	assert.Nil(t, err)
	assert.Equal(t, []byte(DefaultConfig.CipherKey), b)
}
