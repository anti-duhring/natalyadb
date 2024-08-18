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

func init() {
	node1max := HEADER + BTREE_POINTER_SIZE + BTREE_OFFSET_SIZE + 4 + BTREEE_MAX_KEY_SIZE + BTREE_MAX_VAL_SIZE
	if node1max <= BTREE_PAGE_SIZE {
		panic("Node is bigger than the page size")
	}
}
