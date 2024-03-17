package models

type D string

const (
	DATA_TYPE_STRING D = "string"
	DATA_TYPE_INT    D = "int"
)

var DataTypeSizes = map[D]int{
	DATA_TYPE_STRING: 255,
	DATA_TYPE_INT:    11,
}
