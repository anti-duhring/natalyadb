package row_test

import (
	"testing"

	"github.com/anti-duhring/natalyadb/internal/row"
)

func TestSerializeRow(t *testing.T) {
	r := struct {
		Name string
		Age  int
	}{
		Name: "Matt",
		Age:  25,
	}

	serialized := row.SerializeRow(r)

	expected := "0000400040002000Matt25"
	if serialized != expected {
		t.Errorf("Expected %s, got %s", expected, serialized)
	}

	r = struct {
		Name string
		Age  int
	}{
		Name: "Tom Brady",
		Age:  45,
	}

	serialized = row.SerializeRow(r)

	expected = "0000900040002000Tom Brady45"
	if serialized != expected {
		t.Errorf("Expected %s, got %s", expected, serialized)
	}
}
