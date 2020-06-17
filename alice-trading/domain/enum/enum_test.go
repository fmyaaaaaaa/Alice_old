package enum

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// granularity
func TestGranularity_string(t *testing.T) {
	target := S5
	assert.Equal(t, "S5", target.ToString())
}
