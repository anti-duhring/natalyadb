package tests

import (
	"testing"

	"github.com/anti-duhring/natalyadb/internal/index"
)

func TestBTreeInsertAndRetrieve(t *testing.T) {
	// Create a new BTree context
	ctx := index.NewC()

	// Insert a key-value pair
	ctx.Insert("foo", "bar")
	ctx.Insert("tom", "brady")
	ctx.Insert("key", "value")

	assert := NewAssert(t)

	assert.Equal(ctx.Ref["foo"], "bar")
	assert.Equal(ctx.Ref["key"], "value")
	assert.Equal(ctx.Ref["tom"], "brady")
}

func TestBTreeInserAndDelete(t *testing.T) {
	// Create a new BTree context
	ctx := index.NewC()

	// Insert a key-value pair
	ctx.Insert("foo", "bar")
	ctx.Insert("tom", "brady")
	ctx.Insert("key", "value")

	assert := NewAssert(t)

	deleted := ctx.Delete("foo")
	assert.Equal(deleted, true)

	assert.Equal(ctx.Ref["foo"], "")
	assert.Equal(ctx.Ref["key"], "value")
	assert.Equal(ctx.Ref["tom"], "brady")
}
