package strings

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestParsedUrl(t *testing.T) {
	target := "http://test.dummy"
	result := ParsedUrl(target)
	assert.Equal(t, "*url.URL", reflect.TypeOf(result).String())
}
