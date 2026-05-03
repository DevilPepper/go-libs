//go:build prod

package environment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildConstants_prod(t *testing.T) {
	assert.False(t, IS_DEV)
	assert.False(t, IS_QA)
	assert.True(t, IS_PROD)
}
