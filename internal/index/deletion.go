package index

import (
	"bytes"

	. "github.com/anti-duhring/natalyadb/internal/config"
)

// remove a key from a leaf node
func leafDelete(
	new BNode,
	old BNode,
	index uint16,
) {
	new.setHeader(BNODE_LEAF, old.nkeys()-1)
	nodeAppendRange(new, old, 0, 0, index)
	nodeAppendRange(new, old, index, index+1, old.nkeys()-(index+1))
}

func nodeDelete(tree *BTree, node BNode, index uint16, key []byte) BNode {
	// get the child node
	child := node.getPtr(index)
	updated := treeDelete(tree, tree.get(child), key)
	if len(updated.data) == 0 {
		// child doesn't exist
		return BNode{}
	}
	tree.del(child)

	new := BNode{data: make([]byte, BTREE_PAGE_SIZE)}
	// check for merging
	mergeDir, sibling := shouldMerge(tree, node, index, updated)
	if mergeDir == 0 {
		if updated.nkeys() > 0 {
			panic(ERR_INVALID_NODE_SIZE)
		}
		// no merge, just update the child
		nodeReplaceChildLinks(tree, new, node, index, updated)
		return new
	}
	if mergeDir < 0 {
		// left
		merged := BNode{data: make([]byte, BTREE_PAGE_SIZE)}
		nodeMerge(merged, sibling, updated)
		tree.del(node.getPtr(index - 1))
		nodeReplace2Child(new, node, index, tree.new(merged), merged.getKey(0))
		return new
	}
	if mergeDir > 0 {
		// right
		merged := BNode{data: make([]byte, BTREE_PAGE_SIZE)}
		nodeMerge(merged, updated, sibling)
		tree.del(node.getPtr(index + 1))
		nodeReplace2Child(new, node, index, tree.new(merged), sibling.getKey(0))
		return new
	}

	return new
}

// remove a key from tree
func treeDelete(
	tree *BTree,
	node BNode,
	key []byte,
) BNode {
	// find the key
	index := nodeLookup(node, key)
	// act depending on the node type
	switch node.btype() {
	case BNODE_LEAF:
		if !bytes.Equal(node.getKey(index), key) {
			// key doesn't exist
			return BNode{}
		}
		// key exists, delete it
		new := BNode{data: make([]byte, BTREE_PAGE_SIZE)}
		leafDelete(new, node, index)
		return new
	case BNODE_NODE:
		// internal node, delete it from a child node
		return nodeDelete(tree, node, index, key)
	default:
		panic(ERR_INVALID_NODE_TYPE)
	}
}
