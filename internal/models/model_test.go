package models_test

import (
	"fmt"
	"testing"

	"github.com/anti-duhring/natalyadb/internal/models"
)

func TestGetTypes(t *testing.T) {
	m := struct {
		Name string
		Age  int
	}{}

	types := models.GetTypes(m)

	if types["name"] != "string" {
		t.Errorf("Expected string, got %s", types["name"])
	}
	if types["age"] != "int" {
		t.Errorf("Expected int, got %s", types["age"])
	}
}

func TestGetSchema(t *testing.T) {
	a := struct {
		Name string
		Age  int
	}{}

	b := struct {
		Age  int
		Name string
	}{}

	schemaA := models.GetSchema(a)
	schemaB := models.GetSchema(b)

	if fmt.Sprintf("%v", schemaA) != fmt.Sprintf("%v", schemaB) {
		t.Errorf("Expected %v, got %v", schemaA, schemaB)
	}

}
