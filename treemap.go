package tradingdb2

// TreeMapNode - treemap node
type TreeMapNode struct {
	Name     string
	Children map[string]*TreeMapNode
}

// NewTreeMapNode - new TreeMapNode
func NewTreeMapNode(name string) *TreeMapNode {
	return &TreeMapNode{
		Name:     name,
		Children: make(map[string]*TreeMapNode),
	}
}

// AddChild - add a child
func (tree *TreeMapNode) AddChild(name string) (*TreeMapNode, error) {
	_, isok := tree.Children[name]
	if isok {
		return nil, ErrDuplicateChildNode
	}

	child := NewTreeMapNode(name)

	tree.Children[name] = child

	return child, nil
}

// AddChildNode - add a child
func (tree *TreeMapNode) AddChildNode(child *TreeMapNode) error {
	_, isok := tree.Children[child.Name]
	if isok {
		return ErrDuplicateChildNode
	}

	tree.Children[child.Name] = child

	return nil
}

// GetChildEx - add a child
func (tree *TreeMapNode) GetChildEx(name string) *TreeMapNode {
	child, isok := tree.Children[name]
	if isok {
		return child
	}

	child = NewTreeMapNode(name)

	tree.Children[name] = child

	return child
}
