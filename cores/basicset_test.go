package cores

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicSetCore_Push(t *testing.T) {
	setCore, err := BasicSet(nil, nil)

	if err != nil {
		t.Error(err)
	}

	setCore.Push([]byte("test1"))
	setCore.Push([]byte("test1"))
	setCore.Push([]byte("test2"))

	assert.Equal(t, uint64(2), setCore.Count())

}
