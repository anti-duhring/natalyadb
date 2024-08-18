package index

import (
	"unsafe"

	. "github.com/anti-duhring/natalyadb/internal/config"
)

type C struct {
	Tree  BTree
	Ref   map[string]string
	pages map[uint64]BNode
}

func NewC() *C {
	pages := map[uint64]BNode{}
	return &C{
		Tree: BTree{
			get: func(ptr uint64) BNode {
				node, ok := pages[ptr]
				if !ok {
					panic("page not found")
				}
				return node
			},
			new: func(node BNode) uint64 {
				if node.nbytes() > BTREE_PAGE_SIZE {
					panic("node too big")
				}

				key := uint64(uintptr(unsafe.Pointer(&node.data[0])))
				pages[key] = node
				return key
			},
			del: func(ptr uint64) {
				_, ok := pages[ptr]
				if !ok {
					panic("page not found")
				}
				delete(pages, ptr)
			},
		},
		Ref:   map[string]string{},
		pages: pages,
	}
}

func (c *C) Insert(key string, value string) {
	c.Tree.Insert([]byte(key), []byte(value))
	c.Ref[key] = value
}

func (c *C) Delete(key string) bool {
	delete(c.Ref, key)
	return c.Tree.Delete([]byte(key))
}
