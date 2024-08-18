package index

import . "github.com/anti-duhring/natalyadb/internal/config"

// merge 2 nodes into 1
func nodeMerge(new BNode, left BNode, right BNode) {
	new.setHeader(left.btype(), left.nkeys()+right.nkeys())
	nodeAppendRange(new, left, 0, 0, left.nkeys())
	nodeAppendRange(new, right, left.nkeys(), 0, right.nkeys())
}

func shouldMerge(tree *BTree, node BNode, index uint16, updated BNode) (int, BNode) {
	// conditions:
	// 1. if the node is smaller than 1/4 of the page size (arbitrary)
	if updated.nbytes() > BTREE_PAGE_SIZE/4 {
		return 0, BNode{}
	}
	// 2. has sibling and the merged result doesn't exceed the page size
	if index > 0 {
		sibling := tree.get(node.getPtr(index - 1))
		merged := sibling.nbytes() + updated.nbytes() - HEADER
		if merged <= BTREE_PAGE_SIZE {
			return -1, sibling
		}
	}
	if index+1 < node.nkeys() {
		sibling := tree.get(node.getPtr(index + 1))
		merged := updated.nbytes() + sibling.nbytes() - HEADER
		if merged <= BTREE_PAGE_SIZE {
			return 1, sibling
		}
	}

	return 0, BNode{}
}
