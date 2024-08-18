package index

import (
	"errors"
)

var (
	ERR_OUT_OF_RANGE      = errors.New("Index position out of range")
	ERR_INVALID_NODE_TYPE = errors.New("Invalid node type")
	ERR_INVALID_NODE_SIZE = errors.New("Invalid node size")
	ERR_INVALID_KEY       = errors.New("Invalid key")
)
