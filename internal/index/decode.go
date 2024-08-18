package index

import (
	"encoding/binary"

	. "github.com/anti-duhring/natalyadb/internal/config"
)

const (
	BNODE_NODE = 1 // internal nodes without values
	BNODE_LEAF = 2 // leaf nodes with values
)

type BNode struct {
	data []byte
}

// Using little-endien because it allows certain optimizations such as easier incrementing and handling of variable-length data ttypes
// header
func (node BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(node.data)
}

// get the number of keys
func (node BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node.data[2:4])
}

func (node BNode) setHeader(btype uint16, nkeys uint16) {
	binary.LittleEndian.PutUint16(node.data[0:2], btype)
	binary.LittleEndian.PutUint16(node.data[2:4], nkeys)
}

// pointers
func (node BNode) getPtr(index uint16) uint64 {
	if index >= node.nkeys() {
		panic(ERR_OUT_OF_RANGE)
	}

	position := HEADER + BTREE_POINTER_SIZE*index
	return binary.LittleEndian.Uint64(node.data[position:])
}

func (node BNode) setPtr(index uint16, val uint64) {
	if index >= node.nkeys() {
		panic(ERR_OUT_OF_RANGE)
	}
	position := HEADER + BTREE_POINTER_SIZE*index
	binary.LittleEndian.PutUint64(node.data[position:], val)
}

// offset list
func offsetPosition(node BNode, index uint16) uint16 {
	if index < 0 || index > node.nkeys() {
		panic(ERR_OUT_OF_RANGE)
	}

	return HEADER + BTREE_POINTER_SIZE*node.nkeys() + 2*index
}

func (node BNode) getOffset(index uint16) uint16 {
	if index == 0 {
		return 0
	}
	return binary.LittleEndian.Uint16(node.data[offsetPosition(node, index):])
}

func (node BNode) setOffset(index uint16, val uint16) {
	binary.LittleEndian.PutUint16(node.data[offsetPosition(node, index):], val)
}

// key-values
func (node BNode) kvPosition(index uint16) uint16 {
	if 1 > node.nkeys() {
		panic(ERR_OUT_OF_RANGE)
	}

	return HEADER + BTREE_POINTER_SIZE*node.nkeys() + 2*node.nkeys() + node.getOffset(index)
}

func (node BNode) getKey(index uint16) []byte {
	if index >= node.nkeys() {
		panic(ERR_OUT_OF_RANGE)
	}
	position := node.kvPosition(index)
	keyLength := binary.LittleEndian.Uint16(node.data[position:])
	return node.data[position+4:][:keyLength]
}

func (node BNode) getValue(index uint16) []byte {
	if index >= node.nkeys() {
		panic(ERR_OUT_OF_RANGE)
	}
	position := node.kvPosition(index)
	keyLength := binary.LittleEndian.Uint16(node.data[position+0:])
	valueLength := binary.LittleEndian.Uint16(node.data[position+2:])
	return node.data[position+4+keyLength:][:valueLength]
}

// node size in bytes
func (node BNode) nbytes() uint16 {
	return node.kvPosition(node.nkeys())
}
