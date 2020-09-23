package tradingdb2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TreeMapNode(t *testing.T) {
	root := NewTreeMapNode("root")
	assert.NotNil(t, root)

	child0, err := root.AddChild("child0")
	assert.NoError(t, err)
	assert.NotNil(t, child0)

	child1, err := root.AddChild("child1")
	assert.NoError(t, err)
	assert.NotNil(t, child1)

	child0d, err := root.AddChild("child0")
	assert.Error(t, err)
	assert.Nil(t, child0d)

	child2 := NewTreeMapNode("child2")
	assert.NotNil(t, child2)

	child20, err := child2.AddChild("child20")
	assert.NoError(t, err)
	assert.NotNil(t, child20)

	err = root.AddChildNode(child2)
	assert.NoError(t, err)

	err = root.AddChildNode(child2)
	assert.Error(t, err)

	t.Logf("Test_TreeMapNode OK")
}
