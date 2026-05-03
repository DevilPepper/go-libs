//go:build !prod && !qa

package environment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildConstants_dev(t *testing.T) {
	assert.True(t, IS_DEV)
	assert.False(t, IS_QA)
	assert.False(t, IS_PROD)
}
