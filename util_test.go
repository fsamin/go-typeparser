package typeparser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	var strs = List([]string{"peach", "apple", "pear", "plum"})

	assert.Equal(t, 2, strs.Index("pear"))

	assert.False(t, strs.Has("grape"))

	assert.True(t, strs.Any(func(v string) bool {
		return strings.HasPrefix(v, "p")
	}))

	assert.False(t, strs.All(func(v string) bool {
		return strings.HasPrefix(v, "p")
	}))

	assert.EqualValues(t, List([]string{"peach", "apple", "pear"}), strs.Filter(func(v string) bool {
		return strings.Contains(v, "e")
	}))

	assert.EqualValues(t, List([]string{"PEACH", "APPLE", "PEAR", "PLUM"}), strs.Map(strings.ToUpper))
}
