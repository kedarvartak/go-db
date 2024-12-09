package storage

import (
	"fmt"
)

// BTreeNode represents a node in the B-tree
type BTreeNode struct {
	isLeaf   bool         // Whether this is a leaf node
	keys     []int        // Keys stored in this node
	values   []RecordID   // Values (RecordIDs) stored in leaf nodes
	children []*BTreeNode // Child pointers (nil for leaf nodes)
	parent   *BTreeNode   // Parent pointer (nil for root)
}

// BTree represents a B-tree index structure
type BTree struct {
	root   *BTreeNode // Root node of the tree
	degree int        // Minimum degree (minimum number of keys per node)
}

// NewBTree creates a new B-tree with the specified degree
func NewBTree(degree int) *BTree {
	if degree < 2 {
		return nil
	}
	return &BTree{
		root:   nil,
		degree: degree,
	}
}

// Search looks up a key in the B-tree and returns the associated RecordID
func (t *BTree) Search(key int) (*RecordID, bool) {
	if t.root == nil {
		return nil, false
	}
	return t.searchNode(t.root, key)
}

// searchNode recursively searches for a key in a node and its children
func (t *BTree) searchNode(node *BTreeNode, key int) (*RecordID, bool) {
	// Find the first key greater than or equal to target
	i := 0
	for i < len(node.keys) && key > node.keys[i] {
		i++
	}

	// If we found the key in current node
	if i < len(node.keys) && key == node.keys[i] {
		return &node.values[i], true
	}

	// If we haven't found the key and this is a leaf, key doesn't exist
	if node.isLeaf {
		return nil, false
	}

	// Search appropriate child
	// If i == len(node.keys), search rightmost child
	// Otherwise, search child i
	return t.searchNode(node.children[i], key)
}

// Insert adds a new key-value pair to the B-tree
func (t *BTree) Insert(key int, value RecordID) {
	// If tree is empty, create root
	if t.root == nil {
		t.root = &BTreeNode{
			isLeaf: true,
			keys:   make([]int, 0, t.degree),
			values: make([]RecordID, 0, t.degree),
		}
	}

	// If root is full, create new root and split
	if len(t.root.keys) == 2*t.degree-1 {
		newRoot := &BTreeNode{
			isLeaf:   false,
			keys:     make([]int, 0, t.degree),
			children: make([]*BTreeNode, 1, t.degree+1),
		}
		newRoot.children[0] = t.root
		t.root.parent = newRoot
		t.root = newRoot
		t.splitChild(newRoot, 0)
	}

	t.insertNonFull(t.root, key, value)
}

// insertNonFull inserts a key-value pair into a non-full node
func (t *BTree) insertNonFull(node *BTreeNode, key int, value RecordID) {
	if node.isLeaf {
		// Find insertion point
		i := len(node.keys) - 1

		// Make space for new key/value
		node.keys = append(node.keys, 0)
		node.values = append(node.values, RecordID{})

		// Shift elements
		for i >= 0 && node.keys[i] > key {
			node.keys[i+1] = node.keys[i]
			node.values[i+1] = node.values[i]
			i--
		}

		// Insert new key/value
		node.keys[i+1] = key
		node.values[i+1] = value
	} else {
		// Find child to recurse to
		i := len(node.keys) - 1
		for i >= 0 && key < node.keys[i] {
			i--
		}
		i++

		// If child is full, split it
		if len(node.children[i].keys) == 2*t.degree-1 {
			t.splitChild(node, i)
			if key > node.keys[i] {
				i++
			}
		}

		t.insertNonFull(node.children[i], key, value)
	}
}

// splitChild splits a full child node
func (t *BTree) splitChild(parent *BTreeNode, childIndex int) {
	child := parent.children[childIndex]
	medianIdx := t.degree - 1

	// Create new right node
	newRight := &BTreeNode{
		isLeaf:   child.isLeaf,
		keys:     make([]int, 0, t.degree-1),
		values:   make([]RecordID, 0, t.degree-1),
		children: make([]*BTreeNode, 0, t.degree),
		parent:   parent,
	}

	// Move median key and value to parent
	medianKey := child.keys[medianIdx]
	medianValue := child.values[medianIdx]

	// Copy right half to new node
	newRight.keys = append(newRight.keys, child.keys[medianIdx+1:]...)
	newRight.values = append(newRight.values, child.values[medianIdx+1:]...)

	if !child.isLeaf {
		newRight.children = append(newRight.children, child.children[medianIdx+1:]...)
		for _, c := range newRight.children {
			c.parent = newRight
		}
	}

	// Truncate left child
	child.keys = child.keys[:medianIdx]
	child.values = child.values[:medianIdx]
	if !child.isLeaf {
		child.children = child.children[:medianIdx+1]
	}

	// Insert median into parent
	insertIdx := childIndex
	parent.keys = append(parent.keys, 0)
	parent.values = append(parent.values, RecordID{})
	copy(parent.keys[insertIdx+1:], parent.keys[insertIdx:])
	copy(parent.values[insertIdx+1:], parent.values[insertIdx:])
	parent.keys[insertIdx] = medianKey
	parent.values[insertIdx] = medianValue

	// Insert new child into parent
	parent.children = append(parent.children, nil)
	copy(parent.children[childIndex+2:], parent.children[childIndex+1:])
	parent.children[childIndex+1] = newRight
}

// Add this debug helper function
func (t *BTree) String() string {
	if t.root == nil {
		return "Empty tree"
	}
	return t.nodeToString(t.root, 0, "")
}

func (t *BTree) nodeToString(node *BTreeNode, level int, prefix string) string {
	result := fmt.Sprintf("%sLevel %d: Keys: %v Values: %v\n", prefix, level, node.keys, node.values)
	if !node.isLeaf {
		for i, child := range node.children {
			result += fmt.Sprintf("%sChild %d:\n", prefix, i)
			result += t.nodeToString(child, level+1, prefix+"  ")
		}
	}
	return result
}
