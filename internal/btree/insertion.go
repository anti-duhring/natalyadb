package btree

import (
	"bytes"
	"encoding/binary"

	. "github.com/anti-duhring/natalyadb/internal/config"
)

// returns the first child node whose range intersects the key: child[i] <= key
func nodeLookup(node BNode, key []byte) uint16 {
	nkeys := node.nkeys()
	found := uint16(0)
	// the first key is a copy from the parent node, thus it's always less than or equal to the key
	for i := uint16(1); i < nkeys; i++ {
		cmp := bytes.Compare(key, node.getKey(i))
		if cmp <= 0 {
			found = i
		}
		if cmp >= 0 {
			break
		}
	}

	return found
}

// add a new key to a leaf node
func leafInsert(
	new BNode,
	old BNode,
	index uint16,
	key []byte,
	val []byte,
) {
	new.setHeader(BNODE_LEAF, old.nkeys()+1)
	nodeAppendRange(new, old, 0, 0, index)
	nodeAppendKV(new, index, 0, key, val)
	nodeAppendRange(new, old, index+1, index, old.nkeys()-index)
}

// copy multiple key-values into the position
func nodeAppendRange(
	new BNode,
	old BNode,
	dstNew uint16, srcOld uint16, n uint16,
) {
	if srcOld+n <= old.nkeys() || dstNew+n <= new.nkeys() {
		panic(ERR_OUT_OF_RANGE)
	}
	if n == 0 {
		return
	}

	// pointers
	for i := uint16(0); i < n; i++ {
		new.setPtr(dstNew+i, old.getPtr(srcOld+i))
	}
	// offsets
	dstBegin := new.getOffset(dstNew)
	srcBegin := old.getOffset(srcOld)
	// range is [1, n]
	for i := uint16(0); i <= n; i++ {
		offset := dstBegin + old.getOffset(srcOld+i) - srcBegin
		new.setOffset(dstNew+i, offset)
	}
	// key-values
	begin := old.kvPosition(srcOld)
	end := old.kvPosition(srcOld + n)
	copy(new.data[new.kvPosition(dstNew):], old.data[begin:end])
}

// copy a key-value into the position
func nodeAppendKV(
	new BNode,
	index uint16,
	pointer uint64,
	key []byte,
	val []byte,
) {
	// pointers
	new.setPtr(index, pointer)
	// key-values
	position := new.kvPosition(index)
	binary.LittleEndian.PutUint16(new.data[position+0:], uint16(len(key)))
	binary.LittleEndian.PutUint16(new.data[position+2:], uint16(len(val)))
	copy(new.data[position+4:], key)
	copy(new.data[position+4+uint16(len(key)):], val)
	// the offset of the next key
	new.setOffset(index+1, new.getOffset(index)+4+uint16((len(key)+len(val))))
}

// inser a key-value into a node, the result might be split into 2 nodes.
// the caller is responsible for deallocating the input node
// and splitting and allocation result nodes
func treeInsert(
	tree *BTree,
	node BNode,
	key []byte,
	val []byte,
) BNode {
	// the result node. its allowed to be bigger than 1 page and will be split if so
	new := BNode{data: make([]byte, 2*BTREE_PAGE_SIZE)}

	// where to inser the key
	index := nodeLookup(node, key)
	// act depending onn node type
	switch node.btype() {
	case BNODE_LEAF:
		// leaf, node.getKey(index) <= key
		if bytes.Equal(node.getKey(index), key) {
			// key already exists, update the value
			leafInsert(new, node, index, key, val)
		} else {
			// key doesn't exist, insert it
			leafInsert(new, node, index+1, key, val)
		}
	case BNODE_NODE:
		// internal node, insert it to a child node
		nodeInsert(tree, new, node, index, key, val)
	default:
		panic(ERR_INVALID_NODE_TYPE)
	}

	return new
}

func nodeInsert(
	tree *BTree,
	new BNode,
	old BNode,
	index uint16,
	key []byte,
	val []byte,
) {
	// get and deallocate the child node
	childPtr := old.getPtr(index)
	childNode := tree.get(childPtr)
	tree.del(childPtr)
	// recursive insertion to the child node
	childNode = treeInsert(tree, childNode, key, val)
	// split the result
	nsplit, splited := nodeSplit3(childNode)
	// update the child links
	nodeReplaceChildLinks(tree, new, old, index, splited[:nsplit]...)
}

// replace a link with multiple links
func nodeReplaceChildLinks(
	tree *BTree,
	new BNode,
	old BNode,
	index uint16,
	children ...BNode,
) {
	inc := uint16(len(children))
	new.setHeader(BNODE_NODE, old.nkeys()+inc-1)
	nodeAppendRange(new, old, 0, 0, index)
	for i, child := range children {
		nodeAppendKV(new, index+uint16(i), tree.new(child), child.getKey(0), nil)
	}
	nodeAppendRange(new, old, index+inc, index+1, old.nkeys()-(index+1))
}

// replace 2 adjacent links with a single link
func nodeReplace2Child(
	new BNode,
	old BNode,
	index uint16,
	pointer uint64,
	key []byte,
) {
	new.setHeader(BNODE_NODE, old.nkeys()-1)
	nodeAppendRange(new, old, 0, 0, index)
	nodeAppendKV(new, index, pointer, key, nil)
	nodeAppendRange(new, old, index+1, index+2, old.nkeys()-(index+2))
}
