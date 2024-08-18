package btree

import . "github.com/anti-duhring/natalyadb/internal/config"

// split a node if its too big into 1~3 nodes
func nodeSplit3(old BNode) (uint16, [3]BNode) {
	if old.nbytes() <= BTREE_PAGE_SIZE {
		old.data = old.data[:BTREE_PAGE_SIZE]
		return 1, [3]BNode{old}
	}
	left := BNode{make([]byte, 2*BTREE_PAGE_SIZE)} // might be split later
	right := BNode{make([]byte, BTREE_PAGE_SIZE)}
	nodeSplit2(left, right, old)
	if left.nbytes() <= BTREE_PAGE_SIZE {
		left.data = left.data[:BTREE_PAGE_SIZE]
		return 2, [3]BNode{left, right}
	}
	// the left node is still too big, split it again
	leftleft := BNode{make([]byte, BTREE_PAGE_SIZE)}
	middle := BNode{make([]byte, BTREE_PAGE_SIZE)}
	nodeSplit2(leftleft, middle, left)
	if leftleft.nbytes() <= BTREE_PAGE_SIZE {
		panic(ERR_INVALID_NODE_SIZE)
	}

	return 3, [3]BNode{leftleft, middle, right}
}

// split a node into 2 nodes
func nodeSplit2(left BNode, right BNode, old BNode) {
	left.data = old.data[:BTREE_PAGE_SIZE]
	right.data = old.data[BTREE_PAGE_SIZE:]
}
