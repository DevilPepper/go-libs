//go:build qa

package environment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildConstants_qa(t *testing.T) {
	assert.False(t, IS_DEV)
	assert.True(t, IS_QA)
	assert.False(t, IS_PROD)
}
