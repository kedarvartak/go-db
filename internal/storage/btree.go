package storage

type BTreeNode struct {
    isLeaf    bool
    keys      []int
    children  []*BTreeNode
    data      [][]byte
}

type BTree struct {
    root     *BTreeNode
    degree   int
}

func NewBTree(degree int) *BTree {
    return &BTree{
        root:   nil,
        degree: degree,
    }
} 