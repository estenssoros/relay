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
