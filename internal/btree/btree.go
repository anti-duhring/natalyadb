package btree

import . "github.com/anti-duhring/natalyadb/internal/config"

type BTree struct {
	// pointer (a nonzero page number) referencing disk pages
	root uint64
	// callbacks for managing on-disk pages
	get func(uint64) BNode // dereference a pointer to a page
	new func(BNode) uint64 // allocate a new page and return a pointer to it
	del func(uint64)       // deallocate a page
}

func (tree *BTree) Delete(key []byte) bool {
	if len(key) == 0 || len(key) > BTREE_MAX_KEY_SIZE {
		panic(ERR_INVALID_KEY)
	}

	if tree.root == 0 {
		return false
	}

	updated := treeDelete(tree, tree.get(tree.root), key)
	if len(updated.data) == 0 {
		// root doesn't exist
		return false
	}

	tree.del(tree.root)
	if updated.btype() == BNODE_NODE && updated.nkeys() == 1 {
		// the root is no longer needed
		tree.root = updated.getPtr(0)
		return true
	}

	tree.root = tree.new(updated)
	return true
}

func (tree *BTree) Insert(key []byte, val []byte) {
	if len(key) == 0 || len(key) > BTREE_MAX_KEY_SIZE || len(val) > BTREE_MAX_VAL_SIZE {
		panic(ERR_INVALID_KEY)
	}

	if tree.root == 0 {
		// create the first node
		root := BNode{data: make([]byte, BTREE_PAGE_SIZE)}
		root.setHeader(BNODE_LEAF, 2)
		// a dummy key, this makes the tree cover the whole key space thus a lookup will always find a containing node
		nodeAppendKV(root, 0, 0, nil, nil)
		nodeAppendKV(root, 1, 0, key, val)
		tree.root = tree.new(root)
		return
	}

	node := tree.get(tree.root)
	tree.del(tree.root)

	node = treeInsert(tree, node, key, val)
	nsplit, splitted := nodeSplit3(node)

	if nsplit > 1 {
		// the root is split
		root := BNode{data: make([]byte, BTREE_PAGE_SIZE)}
		root.setHeader(BNODE_NODE, nsplit)
		for i, child := range splitted[:nsplit] {
			ptr, key := tree.new(child), child.getKey(0)
			nodeAppendKV(root, uint16(i), ptr, key, nil)
		}
		tree.root = tree.new(root)
		return
	}

	tree.root = tree.new(splitted[0])
}

func init() {
	node1max := HEADER + BTREE_POINTER_SIZE + BTREE_OFFSET_SIZE + 4 + BTREE_MAX_KEY_SIZE + BTREE_MAX_VAL_SIZE
	if node1max <= BTREE_PAGE_SIZE {
		panic("Node is bigger than the page size")
	}
}
