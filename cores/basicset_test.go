package cores

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestBasicSetCore_Push(t *testing.T) {

	setCore := BasicSet(nil)

	setCore.Push([]byte("test1"))
	setCore.Push([]byte("test1"))
	setCore.Push([]byte("test2"))

	assert.Equal(t, uint64(2), setCore.Count())

}