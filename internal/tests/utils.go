package tests

import "testing"

type Assert struct {
	t *testing.T
}

func NewAssert(t *testing.T) *Assert {
	return &Assert{t}
}

func (a *Assert) Equal(expected, actual interface{}) {
	if expected != actual {
		a.t.Fatalf("expected %v, got %v", expected, actual)
	}
}
