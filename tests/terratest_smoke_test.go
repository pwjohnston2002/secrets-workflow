package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerratestEnvironment(t *testing.T) {
	assert.Equal(t, 2, 1+1, "math still works")
}

